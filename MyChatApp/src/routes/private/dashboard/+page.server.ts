import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

const RECEIVER_SERVICE = import.meta.env.VITE_RECEIVER_SERVICE;
const PROTOCOL = import.meta.env.VITE_PROTOCOL;

export const load: PageServerLoad = async ({locals: { supabase } }) => {
    const { data: authData, error: userError } = await supabase.auth.getUser();

    if (userError || !authData) redirect(303, '/auth/signin');

    const { data: loggedUserData} = await supabase.from('User').select('*').eq('user_id', authData.user.id).single();
    const { data: chatData} = await supabase.from('Chat').select('participant1, participant2, blocked')
          .or(`participant1.eq.${loggedUserData.user_id},participant2.eq.${loggedUserData.user_id}`);
    const { data: messageData} = await supabase.from('Message').select('sender, receiver, content, created_at, isRead')
          .or(`sender.eq.${loggedUserData.user_id},receiver.eq.${loggedUserData.user_id}`);

    const ids = [];
    chatData.forEach(chat => {
        const { participant1, participant2 } = chat;
        if (participant1 === loggedUserData.user_id && participant2 === loggedUserData.user_id) ids.push(participant1);
        else ids.push(participant1 === loggedUserData.user_id ? participant2 : participant1);
    });

    const { data: usernames} = await supabase.from('User').select('username, user_id').in('user_id', ids);

    const usernameMap = new Map(usernames.map(user => [user.user_id, user.username]));
    const chats = [];
    if (usernameMap.has(loggedUserData.user_id)) {
      const filteredMessages =  messageData.filter(d => (d.sender === loggedUserData.user_id && d.receiver === loggedUserData.user_id));
      const chat = chatData.filter(c => c.participant1 === loggedUserData.user_id && c.participant2 === loggedUserData.user_id)[0];
      chats.push({username: loggedUserData.username, user_id: loggedUserData.user_id, blocked: false, messages: filteredMessages, unreadCount: 0});
      usernameMap.delete(loggedUserData.user_id);
    }
    for (const chat of chatData) {    
      const filteredMessages =  messageData.filter(d => (d.sender === chat.participant1 && d.receiver === chat.participant2) || (d.sender === chat.participant2 && d.receiver === chat.participant1)).sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
      const unreadCount: number = messageData.filter(msg => !msg.isRead && (msg.sender === chat.participant1 || msg.sender === chat.participant2)).length;
      if (usernameMap.has(chat.participant1)) chats.push({username: usernameMap.get(chat.participant1), user_id: chat.participant1, blocked: chat.blocked, messages: filteredMessages, unreadCount: unreadCount});
      else if (usernameMap.has(chat.participant2)) chats.push({username: usernameMap.get(chat.participant2), user_id: chat.participant2, blocked: chat.blocked, messages: filteredMessages, unreadCount: unreadCount});
    }
	  return { user: loggedUserData, chats: chats};
};

export const actions = {
  getUserId: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const username = formData.get('username') as string;
    if (typeof username !== 'string' || !username.trim()) {
      return null;
    }
    const { data, error } = await supabase.from('User').select('user_id').eq('username', username).single();
    if (error || !data) return null;
    return data.user_id;
  },

  createChat: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const p1 = formData.get('participant1') as string;
    const p2 = formData.get('participant2') as string;
    const { data, error } = await supabase
    .from('Chat').select('participant1, participant2')
    .or(`and(participant1.eq.${p1},participant2.eq.${p2}),` + `and(participant1.eq.${p2},participant2.eq.${p1})`)
    .single();
    if (error || !data) {
      await supabase.from('Chat').insert([{participant1: p1, participant2: p2}]).single();
    }    
  },

  logout: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const user_id = formData.get('user_id');

  
    if (user_id) {
      try {
        const response = await fetch(PROTOCOL + RECEIVER_SERVICE + "/close", {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ user_id: user_id }),
        });
  
        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }
      } catch (error) {
        console.log('Error during fetch:', error);
      } 
    } else {
      console.log('No user_id provided.');
    }
    await supabase.auth.signOut();
  }
};

