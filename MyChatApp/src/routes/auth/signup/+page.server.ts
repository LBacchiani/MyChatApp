import { redirect } from '@sveltejs/kit'

import type { Actions } from './$types.js'

export const actions: Actions = {
  signup: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData()
    const email = formData.get('email') as string
    const password = formData.get('password') as string
    const username = formData.get('username') as string

    const { data, emailError } = await supabase.auth.signUp({ email, password })
    const { usernameError } = await supabase.from('User').insert({email: email, username: username, user_id: data.user.id})
    if (emailError || usernameError) {
      redirect(303, '/auth/error')
    } else {
      redirect(303, '/private/dashboard')
    }
  }
}