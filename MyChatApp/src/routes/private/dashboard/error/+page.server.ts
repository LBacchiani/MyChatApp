import type { PageServerLoad } from './$types.js';

export const load: PageServerLoad = async ({locals: { supabase } }) => {
    await supabase.auth.signOut();
};

