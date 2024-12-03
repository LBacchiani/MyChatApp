<script lang="ts">
  import { goto } from '$app/navigation';

  async function handleKeydown(event: KeyboardEvent): Promise<void> {
    if (event.key === 'Enter') {
      event.preventDefault(); 
      const formData = new FormData();
      formData.append('username', username2);

      const response = await fetch('?/getUserId', {
        method: 'POST',
        body: formData,
      });
      console.log(response)
      const data = await response.json();
      const user_id = JSON.parse(data.data)[0];
      if (!user_id) {
        alert('Username not found');
        return;
      }
      const isDuplicate = reactiveChats.some(item => item.username === username2 && item.user_id === user_id);
      if (!isDuplicate) {
        supabase.from('Chat').insert([{participant1: user.user_id, participant2: user_id}]);
        reactiveChats = [...reactiveChats, {username: username2, user_id: username2}];
      }
    }
  }

  export let data;
  let username2: string;
  let reactiveChats: any[] = [];
  $: ({supabase, user, chats} = data);
  $: reactiveChats = chats;
  $: logout = async () => {
    supabase.auth.signOut();
    goto('/auth/signin')
  };

</script>
  
<div class="flex h-screen bg-gray-100">
  <!-- Sidebar -->
  <aside class="w-64 bg-blue-600 text-white flex flex-col">
    <div class="flex items-center justify-center h-16 border-b border-blue-500">
      <h2 class="text-2xl font-bold">MyChatApp</h2>
    </div>
    <div class="p-4">
      <!-- Search Bar for Users -->
      <input
        type="text"
        id="searchInput"
        placeholder="Search users to chat"
        class="w-full p-2 rounded-lg text-gray-700"
        bind:value={username2}
        on:keydown={handleKeydown}
      />
    </div>
    <nav class="flex-grow overflow-y-auto">
      <ul class="space-y-2 p-4">
        {#each reactiveChats as chat}
          <li class="flex items-center p-2 bg-blue-500 rounded-lg hover:bg-blue-400">
            <img
              src="../../../user-icon.svg" 
              alt="Avatar"
              class="w-10 h-10 rounded-full border border-gray-300 mr-3"
            />
            <span class="text-white">{chat.username}</span>
          </li>
        {/each}
      </ul>
    </nav>
    <div class="border-t border-blue-500 p-4">
      <button
        class="w-full bg-gray-700 py-2 px-4 text-white rounded-lg hover:bg-gray-600 mb-2">
        Settings
      </button>
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
      <h1 class="text-2xl font-bold text-gray-800">Welcome, {user.username}!</h1>
      <div class="flex items-center space-x-4">
        <span class="text-sm text-gray-500">Today's Date: {new Date().toLocaleDateString()}</span>
        <img
          src="../../../user-icon.svg"
          alt="User Avatar"
          class="w-10 h-10 rounded-full border border-gray-300"
        />
      </div>
    </header>

    <section class="mt-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Upcoming Meetings Card -->
        <div class="bg-white p-4 rounded-lg shadow">
          <h2 class="text-xl font-bold text-blue-600">Upcoming Meetings</h2>
          <p class="text-gray-600 mt-2">You have no meetings scheduled for today.</p>
        </div>
      </div>
    </section>
  </main>
</div>
