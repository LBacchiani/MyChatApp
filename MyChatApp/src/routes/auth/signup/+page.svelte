<script>
  import { page } from '$app/stores';

  let errorMessage = '';
  $: errorMessage = $page.url.searchParams.get('error') || '';
</script>

<div class="flex flex-col items-center justify-center h-screen bg-[#E5F9F2] p-6">
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
    <h1 class="text-4xl font-bold text-[#25d366] mb-6">Register for MyChatApp</h1>

    {#if errorMessage}
      <div class="bg-[#ffefee] text-[#e60000] p-4 rounded-lg mb-6">
        {errorMessage}
      </div>
    {/if}

    <form method="POST" action="?/signup" class="space-y-6">
      <div>
        <label for="username" class="block text-left text-lg text-gray-700">Username</label>
        <input
          id="username"
          name="username"
          type="text"
          required
          class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#25d366]"
        />
      </div>
  
      <div>
        <label for="email" class="block text-left text-lg text-gray-700">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          required
          class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#25d366]"
        />
      </div>
  
      <div>
        <label for="password" class="block text-left text-lg text-gray-700">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          required
          class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#25d366]"
        />
      </div>
  
      <div class="space-y-4">
        <button 
          type="submit"
          class="w-full bg-[#25d366] text-white py-2 rounded-lg hover:bg-[#128C7E] focus:outline-none">
          Register
        </button>
        <p class="text-sm text-gray-500">
          Already have an account? <a href="/auth/signin" class="text-[#25d366] hover:underline">Login</a>
        </p>
      </div>
    </form>
  </div>

  <footer class="mt-6 text-sm text-gray-500">
    <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
  </footer>
</div>
