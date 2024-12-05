import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({locals: { supabase } }) => {

    const { data: authData, error: userError } = await supabase.auth.getUser();

    if (userError || !authData) redirect(303, '/auth/signin');

    const { data: loggedUserData} = await supabase.from('User').select('*').eq('user_id', authData.user.id).single();
    const { data: chatData} = await supabase.from('Chat').select('id, participant1, participant2, blocked')
          .or(`participant1.eq.${loggedUserData.user_id},participant2.eq.${loggedUserData.user_id}`);
    const { data: messageData} = await supabase.from('Message').select('sender, receiver, content, chat_id, created_at')
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
      const chat = chatData.filter(c => c.participant1 === loggedUserData.user_id && c.participant2 === loggedUserData.user_id)[0];
      chats.push({username: loggedUserData.username, user_id: loggedUserData.user_id, chat_id: chat.id, blocked: false, messages: messageData.filter(d => d.chat_id === chat.id)});
      usernameMap.delete(loggedUserData.user_id);
    }
    for (let chat of chatData) {    
      if (usernameMap.has(chat.participant1)) chats.push({username: usernameMap.get(chat.participant1), user_id: chat.participant1, chat_id: chat.id, blocked: chat.blocked, messages: messageData.filter(d => d.chat_id === chat.id)});
      else if (usernameMap.has(chat.participant2)) chats.push({username: usernameMap.get(chat.participant2), user_id: chat.participant1, chat_id: chat.id, blocked: chat.blocked, messages: messageData.filter(d => d.chat_id === chat.id)});
    }
    console.log(chats[1])
	  return { user: loggedUserData, chats: chats};
};

export const actions = {
  getUserId: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const username = formData.get('username') as string;
    if (typeof username !== 'string' || !username.trim()) {
      return null;
    }
    const { data, error } = await supabase
      .from('User')
      .select('user_id')
      .eq('username', username)
      .single();
    if (error || !data) {
      return null;
    }
    return data.user_id;
  },

  createChat: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const p1 = formData.get('participant1') as string;
    const p2 = formData.get('participant2') as string;
    await supabase.from('Chat').insert([{participant1: p1, participant2: p2}]);
  },

  logout: async ({locals: { supabase } }) => {
    supabase.auth.signOut();
  }, 
};