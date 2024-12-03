import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { json, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({locals: { supabase } }) => {

    const { data: authData, error: userError } = await supabase.auth.getUser();

    if (userError || !authData) redirect(303, '/auth/signin');

    const { data: loggedUserData} = await supabase
        .from('User')
        .select('*')
        .eq('user_id', authData.user.id)
        .single();

    const { data: chatData} = await supabase
        .from('Chat')
        .select('participant1, participant2, blocked')
        .or(`participant1.eq.${loggedUserData.user_id},participant2.eq.${loggedUserData.user_id}`);
    const chats = [];

    for (let chat of chatData) {
        const { data: userData} = await supabase
            .from('User')
            .select('username, user_id')
            .neq('user_id', authData.user.id)
            .or(`user_id.eq.${chat.participant1},user_id.eq.${chat.participant2}`);
        chats.push(userData[0])
    }
	return { user: loggedUserData, chats: chats};
};

export const actions = {
  getUserId: async ({ request, locals: { supabase } }) => {
    // Get form data
    const formData = await request.formData();
    const username = formData.get('username') as string;

    // Validate username
    if (typeof username !== 'string' || !username.trim()) {
      return null;
    }

    // Query the Supabase database
    const { data, error } = await supabase
      .from('User')
      .select('user_id')
      .eq('username', username)
      .single();

    // Handle Supabase error or no data found
    if (error || !data) {
      return null;
    }

    // Return only serializable data (a plain object)
    return data.user_id;
  },
};