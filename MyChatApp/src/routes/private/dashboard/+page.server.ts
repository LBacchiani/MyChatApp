import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({locals: { supabase } }) => {

    const {data, error} = await supabase.auth.getUser();
    if (error) {
        redirect(303, '/auth/signin');
    }
	return { user: data.user };
};