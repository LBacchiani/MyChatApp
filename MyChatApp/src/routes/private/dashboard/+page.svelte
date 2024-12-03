<script lang="ts">
  import { goto } from '$app/navigation';

  async function fetchUsername(user_id: string) {
    return await supabase
        .from('User')
        .select('username')
        .eq('user_id', user_id)
  }

	export let data;
  let username: string;
	$: ({ supabase,user } = data);
  $: fetchUsername(user.id).then(response => username = response.data![0].username);
  




	$: logout = async () => {
		const { error } = await supabase.auth.signOut();
		if (error) {
			console.error(error);
		}
    goto('/auth/signin')
	};
</script>
  
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-blue-600 text-white flex flex-col">
      <div class="flex items-center justify-center h-16 border-b border-blue-500">
        <h2 class="text-2xl font-bold">MyChatApp</h2>
      </div>
      <nav class="flex-grow">
        <ul class="space-y-2 p-4">
          <li>
            <a href="/dashboard" class="flex items-center px-4 py-2 rounded-lg hover:bg-blue-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h11M9 21H3M21 10H15m6 0a9 9 0 11-9-9 9 9 0 019 9z" />
              </svg>
              Dashboard
            </a>
          </li>
          <li>
            <a href="/chat" class="flex items-center px-4 py-2 rounded-lg hover:bg-blue-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v7a2 2 0 01-2 2h-2m-4 0h-4m4 0V8m0 0h4M9 21V8a2 2 0 012-2h2m-4 2H5a2 2 0 00-2 2v7a2 2 0 002 2h4" />
              </svg>
              Chats
            </a>
          </li>
          <li>
            <a href="/settings" class="flex items-center px-4 py-2 rounded-lg hover:bg-blue-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m-6-8h6m-4-4h4m0 12h4M2 20h2m2 0h2m2 0h2m2 0h2m2 0h2m2 0h2" />
              </svg>
              Settings
            </a>
          </li>
        </ul>
      </nav>
      <div class="border-t border-blue-500 p-4">
        <button
          on:click={logout}
          class="w-full bg-red-600 py-2 px-4 text-white rounded-lg hover:bg-red-500">
          Logout
        </button>
      </div>
    </aside>
  
    <!-- Main Content -->
    <main class="flex-grow p-6">
      <header class="flex items-center justify-between bg-white p-4 rounded-lg shadow">
        <h1 class="text-2xl font-bold text-gray-800">Welcome, {username}!</h1>
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-500">Today's Date: {new Date().toLocaleDateString()}</span>
          <img
            src="https://via.placeholder.com/40"
            alt="User Avatar"
            class="w-10 h-10 rounded-full border border-gray-300"
          />
        </div>
      </header>
  
      <section class="mt-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            
          <div class="bg-white p-4 rounded-lg shadow">
            <h2 class="text-xl font-bold text-blue-600">Your Activity</h2>
            <p class="text-gray-600 mt-2">Check out your latest messages and notifications here.</p>
          </div>
  
          <!-- Another Card -->
          <div class="bg-white p-4 rounded-lg shadow">
            <h2 class="text-xl font-bold text-blue-600">Upcoming Meetings</h2>
            <p class="text-gray-600 mt-2">You have no meetings scheduled for today.</p>
          </div>
  
          <!-- Another Card -->
          <div class="bg-white p-4 rounded-lg shadow">
            <h2 class="text-xl font-bold text-blue-600">Explore Features</h2>
            <p class="text-gray-600 mt-2">Discover more tools to help you connect with your team.</p>
          </div>
        </div>
      </section>
    </main>
  </div>
  
  