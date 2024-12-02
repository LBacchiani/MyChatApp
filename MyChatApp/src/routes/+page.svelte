<script> 
  import { goto } from '$app/navigation';
</script>
<div class="flex flex-col items-center justify-center h-screen bg-blue-100 p-6">
    <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
      <h1 class="text-4xl font-bold text-blue-600 mb-6">MyChatApp</h1>
  
      <p class="text-lg text-gray-700 mb-6">
        Welcome to MyChatApp! The best way to chat with friends.
      </p>
      <div class="space-y-4">
        <button 
          class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-500 focus:outline-none"
          on:click={() => goto('/register')}>
          Register
        </button>
        <button 
          class="w-full bg-gray-600 text-white py-2 rounded-lg hover:bg-gray-500 focus:outline-none"
          on:click={() => goto('/login')}>
          Login
        </button>
      </div>
    </div>
    <footer class="mt-6 text-sm text-gray-500">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </div>
  