<div class="flex flex-col items-center justify-center h-screen bg-blue-100 p-6">
    <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
      <h1 class="text-4xl font-bold text-blue-600 mb-6">Login to MyChatApp</h1>  
      <form method="POST" action="?/signin" class="space-y-6">
        <div>
          <label for="email" class="block text-left text-lg text-gray-700">Email</label>
          <input
            id="email"
            name="email"
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
            Don't have an account? <a href="/auth/signup" class="text-blue-600 hover:underline">Register</a>
          </p>
        </div>
      </form>
    </div>
    <footer class="mt-6 text-sm text-gray-500">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </div>
