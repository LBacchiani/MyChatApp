import { redirect } from '@sveltejs/kit'

import type { Actions } from './$types.js'
import Redis from 'ioredis';


export const actions: Actions = {
  signin: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData()
    const email = formData.get('email') as string
    const password = formData.get('password') as string

    const { data, error } = await supabase.auth.signInWithPassword({ email, password })

    const user = data.user;
    if (user) {
      const redis = new Redis({ host: '127.0.0.1', port: 6379 });
      redis.del(user.id).then(res => {
        console.log(res);
        redis.quit();
      });
    }

    if (error) {
      console.error(error)
      redirect(303, `/auth/signin?error=${encodeURIComponent('Wrong username or password')}`)
    } else {
      redirect(303, '/private/dashboard')
    }
  }
}