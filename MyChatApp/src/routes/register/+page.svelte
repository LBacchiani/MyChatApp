<script>
  import { goto } from '$app/navigation';   
  import {supabase} from '$lib/supabase';
    let email = '';
    let password = '';
    let username = '';
    let errorMessage = '';
  
    // Handle form submission
    const handleSubmit = async (event) => {
  event.preventDefault(); 

  if (!email || !password || !username) {
    errorMessage = 'All fields are required.';
    return;
  }

  const { user, error: signUpError } = await supabase.auth.signUp({
    email: email,
    password: password,
  });

  if (signUpError) {
    errorMessage = signUpError.message || 'Something went wrong during registration.';
    return;
  }

  const { insertError } = await supabase.from('User').insert({username: username, email: email });
  if (insertError) {
    errorMessage = 'Username already exists';
    return;
  }

  goto('/login');
};

  </script>
  
  <div class="flex flex-col items-center justify-center h-screen bg-blue-100 p-6">
    <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
      <h1 class="text-4xl font-bold text-blue-600 mb-6">Register for MyChatApp</h1>
  
      {#if errorMessage}
        <div class="bg-red-100 text-red-700 p-4 rounded-lg mb-4">
          <strong>{errorMessage}</strong>
        </div>
      {/if}
  
      <form on:submit={handleSubmit} class="space-y-6">
        <div>
          <label for="username" class="block text-left text-lg text-gray-700">Username</label>
          <input
            id="username"
            type="text"
            bind:value={username}
            required
            class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
  
        <div>
          <label for="email" class="block text-left text-lg text-gray-700">Email</label>
          <input
            id="email"
            type="email"
            bind:value={email}
            required
            class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
  
        <div>
          <label for="password" class="block text-left text-lg text-gray-700">Password</label>
          <input
            id="password"
            type="password"
            bind:value={password}
            required
            class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
  
        <div class="space-y-4">
          <button 
            type="submit"
            class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-500 focus:outline-none">
            Register
          </button>
          <p class="text-sm text-gray-500">
            Already have an account? <a href="/login" class="text-blue-600 hover:underline">Login</a>
          </p>
        </div>
      </form>
    </div>
  
    <footer class="mt-6 text-sm text-gray-500">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </div>