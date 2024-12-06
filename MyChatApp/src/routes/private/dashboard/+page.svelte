<script lang="ts">
  const SENDER_SERVICE = import.meta.env.VITE_SENDER_SERVICE;
  export let data;

  function formatDate(dateString: string) {
    const date = new Date(dateString);
    return  date.getDate() + "/" + date.getMonth() + "/" + date.getFullYear() + " | " + date.getHours() + ":" + date.getMinutes();
  }

  async function sendMessage() {
    if (!newMessage.trim()) return; // Prevent sending empty messages

    try {
      const response = await fetch(SENDER_SERVICE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({sender: user.user_id, receiver: selectedChat.user_id, content: newMessage,}),
      });

      const result = await response.json();

      if (result.success) {
        selectedChat.messages = [...selectedChat.messages, { 
          sender: user.user_id,
          content: newMessage,
          created_at: new Date().toISOString() 
        }];
        newMessage = '';
      } else {
        alert('Message sending failed');
      }
    } catch (error: any) {
      alert('Error sending message: ' + error.message);
    }
  }

  async function handleKeydown(event: KeyboardEvent): Promise<void> {
    if (event.key === 'Enter') {
      event.preventDefault(); 
      const formData = new FormData();
      formData.append('username', username2);

      const response = await fetch('?/getUserId', {
        method: 'POST',
        body: formData,
      });
      const data = await response.json();
      const user_id = JSON.parse(data.data)[0];
      if (!user_id) {
        alert('Username not found');
        return;
      }
      const isDuplicate = chats.some(item => item.user_id === user_id);
      
      if (!isDuplicate) {
        const formData = new FormData();
        formData.append('participant1', user.user_id);
        formData.append('participant2', user_id);
        const response = await fetch('?/createChat', {
          method: 'POST',
          body: formData,
        });   
        const result = await response.json();
        chats = [...chats, {username: username2, user_id: user_id, blocked: false, messages: []}];
        username2 = ''
      }
    }
  }
  let username2: string;
  let newMessage = '';
  const selectChat = (chat) => selectedChat = chat;
  $: user = data.user;
  $: chats = data.chats;
  $: selectedChat = chats[0];
  let sidebarOpen: boolean = false;
</script>

<div class="flex h-screen bg-gray-100">
  <!-- Sidebar (Mobile) -->
  <aside
    class="lg:w-64 w-64 lg:translate-x-0 bg-blue-600 text-white flex flex-col fixed inset-y-0 left-0 transform transition-transform duration-300 ease-in-out"
    class:translate-x-[-100%]="{!sidebarOpen}"
    class:translate-x-0="{sidebarOpen}"
  >
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
        {#each chats as chat}
          <button
            class="flex items-center p-2 w-full text-left bg-blue-500 rounded-lg hover:bg-blue-400 cursor-pointer"
            on:click={() => selectChat(chat)}
          >
            <img
              src="../../../user-icon.svg"
              alt="Avatar"
              class="w-10 h-10 rounded-full border border-gray-300 mr-3"
            />
            <span class="text-white">{chat.username}</span>
          </button>
        {/each}
      </ul>
    </nav>
    <div class="border-t border-blue-500 p-4">
      <button class="w-full bg-gray-700 py-2 px-4 text-white rounded-lg hover:bg-gray-600 mb-2">Settings</button>
      <form method="post" action="?/logout">
        <button
          type="submit"
          class="w-full bg-red-600 py-2 px-4 text-white rounded-lg hover:bg-red-500"
        >
          Logout
        </button>
      </form>
    </div>
  </aside>

  <!-- Hamburger Menu for Mobile -->
  <button
    class="lg:hidden fixed top-4 left-4 z-50 bg-blue-600 text-white p-2 rounded-lg"
    on:click={() => (sidebarOpen = !sidebarOpen)}
  >
    ☰
  </button>

  <!-- Main Content -->
  <main class="flex-1 flex flex-col h-screen overflow-hidden p-4 lg:p-6 ml-0 lg:ml-64">
    <header class="flex items-center justify-between bg-white p-4 rounded-lg shadow shrink-0">
      <h1 class="text-xl lg:text-2xl font-bold text-gray-800">Welcome, {user.username}!</h1>
      <div class="flex items-center space-x-4">
        <span class="hidden md:inline text-sm text-gray-500">Today's Date: {new Date().toLocaleDateString()}</span>
        <img
          src="../../../user-icon.svg"
          alt="User Avatar"
          class="w-8 h-8 lg:w-10 lg:h-10 rounded-full border border-gray-300"
        />
      </div>
    </header>

    <!-- Chat Conversation Section -->
    <section class="mt-4 lg:mt-6 flex-1 min-h-0 flex flex-col">
      {#if selectedChat}
        <div class="bg-white p-4 rounded-lg shadow-lg flex flex-col h-full">
          <h2 class="text-lg lg:text-xl font-bold text-blue-600 mb-4">{selectedChat.username}</h2>
          <div class="flex-1 min-h-0 overflow-y-auto space-y-4">
            {#each selectedChat.messages as message}
              <div class="flex {message.sender === user.user_id ? 'justify-end' : 'justify-start'}">
                <div class="max-w-[75%] p-3 rounded-lg {message.sender === user.user_id ? 'bg-blue-600 text-white' : 'bg-gray-200 text-black'}">
                  <p class="text-sm lg:text-base">{message.content}</p>
                  <span class="text-xs lg:text-sm {message.sender === user.user_id ? 'text-white/80' : 'text-black/60'}">{formatDate(message.created_at)}</span>
                </div>
              </div>
            {/each}
          </div>
          <div class="mt-4 flex items-center">
            <input
              type="text"
              class="flex-1 p-2 rounded-lg border border-gray-300"
              placeholder="Type a message..."
              bind:value={newMessage}
            />
            <button
              class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-500 ml-2"
              on:click={sendMessage}
            >
              Send
            </button>
          </div>
        </div>
      {:else}
        <div class="bg-white p-4 rounded-lg shadow h-full flex items-center justify-center">
          <h2 class="text-lg lg:text-xl font-bold text-blue-600">Select a chat to start messaging</h2>
        </div>
      {/if}
    </section>

    <footer class="mt-4 lg:mt-6 text-sm text-center text-gray-600 shrink-0">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </main>
</div>
