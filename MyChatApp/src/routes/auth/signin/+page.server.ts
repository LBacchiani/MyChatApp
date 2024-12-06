import { redirect } from '@sveltejs/kit'

import type { Actions } from './$types.js'

export const actions: Actions = {
  signin: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData()
    const email = formData.get('email') as string
    const password = formData.get('password') as string

    const { error } = await supabase.auth.signInWithPassword({ email, password })
    if (error) {
      console.error(error)
      redirect(303, `/auth/signin?error=${encodeURIComponent('Wrong username or password')}`)
    } else {
      redirect(303, '/private/dashboard')
    }
  }
}