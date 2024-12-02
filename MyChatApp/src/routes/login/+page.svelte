<script>
  import { goto } from '$app/navigation';
    
    // Handle form submission
    const handleLogin = (event) => {
      event.preventDefault(); // Prevent default form submission
      
      const username = event.target.username.value;
      const password = event.target.password.value;
  
      // For now, just log the credentials to the console
      console.log('Username:', username);
      console.log('Password:', password);
  
      // Simulate successful login and navigate to home or dashboard
      // For now, you can navigate to a different page on success
      goto('/');
    };
  </script>
  
  <div class="flex flex-col items-center justify-center h-screen bg-blue-100 p-6">
    <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
      <h1 class="text-4xl font-bold text-blue-600 mb-6">Login to MyChatApp</h1>
  
      <form on:submit={handleLogin} class="space-y-6">
        <div>
          <label for="username" class="block text-left text-lg text-gray-700">Username</label>
          <input
            id="username"
            name="username"
            type="text"
            required
            class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label for="password" class="block text-left text-lg text-gray-700">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            class="w-full p-3 mt-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
  
        <div class="space-y-4">
          <button 
            type="submit"
            class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-500 focus:outline-none">
            Login
          </button>
          <p class="text-sm text-gray-500">
            Don't have an account? <a href="/register" class="text-blue-600 hover:underline">Register</a>
          </p>
        </div>
      </form>
    </div>
    <footer class="mt-6 text-sm text-gray-500">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </div>
