<script lang="ts">
	import { onMount, tick } from 'svelte';

  const SENDER_SERVICE = import.meta.env.VITE_SENDER_SERVICE;
  const RECEIVER_SERVICE = import.meta.env.VITE_RECEIVER_SERVICE;
  const PROTOCOL = import.meta.env.VITE_PROTOCOL;
  export let data;

  function formatDate(dateString: string) {
    const date = new Date(dateString);
    return  date.getDate() + "/" + date.getMonth() + "/" + date.getFullYear() + " | " + date.getHours() + ":" + (date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes());
  }

  async function scrollToBottom(chatContainer: HTMLDivElement | null): Promise<void> {
    if (chatContainer) {
      await tick();
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }
  }

  function createSocket(user_id: string) {
    let socket: WebSocket | null = new WebSocket(`ws://${RECEIVER_SERVICE}/connect?user_id=${user_id}`); 
    socket.onmessage = function(event) {
        const payload = JSON.parse(event.data);
        const message = JSON.parse(payload);
        console.log(event.data);
        const chat = chats.filter(chat => chat.user_id === message.sender)[0];
        let messages: any;
        if (message.type === 'ack') messages = chat.messages.map(message => ({...message, isRead: true }));
        else {
          messages = [...chat.messages, { sender: message.sender, content: message.content, created_at: new Date().toISOString(), isRead: false }];
          chat.unreadCount = chat.unreadCount + 1;
        }
        chat.messages = messages;
        chats = [...chats]
    };

    socket.onerror = async function(event) {
        if (socket) {
          socket.close();
          socket = null;
        }
        console.log("an error occurred, trying to reconnect....");
        setTimeout(() => createSocket(user_id), 5000);
    };
  }


  async function sendMessage() {
    if (!newMessage.trim()) return; // Prevent sending empty messages
    try {
      const msg = { sender: user.user_id, receiver: chats[selectedChat].user_id, content: newMessage, isRead: user.user_id === chats[selectedChat].user_id }
      const response = await fetch(PROTOCOL + SENDER_SERVICE, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(msg) });
      const result = await response.json();
      if (result.success) {
        const msg = { sender: user.user_id, content: newMessage, created_at: new Date().toISOString(), isRead: user.user_id === chats[selectedChat].user_id};
        chats[selectedChat].messages = [...chats[selectedChat].messages, msg];
        chats = [...chats];
        newMessage = '';
        scrollToBottom(chatContainer);
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
      const response = await fetch('?/getUserId', { method: 'POST', body: formData });
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
        chats = [...chats, { username: username2, user_id: user_id, blocked: false, messages: [] }];
        username2 = ''
      }
    }
  }
  //////////////////////////////
  let username2: string;
  let newMessage = '';
  let sidebarOpen: boolean;
  let mobile: boolean;
  let selectedChat: number = -1;
  let chatContainer: HTMLDivElement | null = null;
  const selectChat = (index: number) => selectedChat = index;
  $: user = data.user;
  $: chats = data.chats;
  $: selectedChat, scrollToBottom(chatContainer);
  $: {
    if (selectedChat > -1) {
      if (chats[selectedChat].unreadCount > 0) {
        fetch(PROTOCOL + SENDER_SERVICE, {
          method: 'POST', headers: {'Content-Type': 'application/json',},
          body: JSON.stringify({ receiver: chats[selectedChat].user_id, sender: user.user_id, type: 'ack' }),
        });
        chats[selectedChat].unreadCount = 0;
      }
    }
  }
  onMount(() => {
    mobile = /Mobi|Android|iPhone|iPad|iPod/.test(navigator.userAgent);
    sidebarOpen = !mobile;
    createSocket(user.user_id);
    scrollToBottom(chatContainer);
    if (chats.length > 0) selectedChat = 0;
  });
</script>

<div class="flex h-screen bg-gray-100">
  <!-- Sidebar -->
  <aside
    class="lg:w-64 min-w-[16rem] bg-[#25d366] text-white flex flex-col fixed inset-y-0 left-0 transform transition-transform duration-300 ease-in-out"
    class:translate-x-[-100%]={!sidebarOpen}
  >
    <div class="flex items-center justify-between p-4 border-b border-[#128C7E]">
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
        {#each chats as chat, index}
          <button 
            class="flex items-center p-2 w-full text-left bg-[rgb(131,248,174)] rounded-lg hover:bg-[#128C7E] cursor-pointer"
            on:click={() => selectChat(index)}
            aria-label="Select chat with {chat.username}"
            aria-pressed={selectedChat === index ? 'true' : 'false'}
          >
            <img
              src="../../../user-icon.svg" 
              alt="Avatar"
              class="w-10 h-10 rounded-full border border-gray-300 mr-3"
            />
            <span class="text-white">{chat.username}</span>
            {#if chat.unreadCount > 0}
            <span
              class="bg-red-500 text-white text-xs font-bold rounded-full px-2 py-1">
              {chat.unreadCount}
            </span>
          {/if}
          </button>
        {/each}
      </ul>
    </nav>
    <div class="border-t border-[#128C7E] p-4">
      <button
        class="w-full bg-gray-700 py-2 px-4 text-white rounded-lg hover:bg-gray-600 mb-2">
        Settings
      </button>
      <form method="post" action="?/logout">
        <input type="hidden" name="user_id" value={user.user_id}>
        <button
          type="submit"
          class="w-full bg-red-600 py-2 px-4 text-white rounded-lg hover:bg-red-500">
          Logout
        </button>
      </form>
    </div>
  </aside>

  {#if mobile}
    <button
      class="lg:hidden fixed top-4 left-4 z-50 text-{sidebarOpen ? 'white' : 'green-500'} text-3xl"
      on:click={() => sidebarOpen = !sidebarOpen}
      aria-label="Toggle sidebar"
    >
      <span class="sr-only">Toggle sidebar</span>
      ☰
    </button>
  {/if}

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

    <section class="mt-4 lg:mt-6 shrink-0">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 lg:gap-6">
        <!-- Upcoming Meetings Card -->
        <div class="bg-white p-4 rounded-lg shadow">
          <h2 class="text-lg lg:text-xl font-bold text-[#25d366]">Upcoming Meetings</h2>
          <p class="text-gray-600 mt-2">You have no meetings scheduled for today.</p>
        </div>
      </div>
    </section>

    <!-- Chat Conversation Section -->
    <section class="mt-4 lg:mt-6 flex-1 min-h-0 flex flex-col">
      {#if chats[selectedChat]}
        <!-- Conversation Section -->
        <div class="bg-white p-4 rounded-lg shadow-lg flex flex-col h-full">
          <h2 class="text-lg lg:text-xl font-bold text-[#25d366] mb-4">{chats[selectedChat].username}</h2>
          <div class="flex-1 min-h-0 overflow-y-auto space-y-4" bind:this={chatContainer}>
            <!-- Display Messages -->
            {#each chats[selectedChat].messages as message}
              <div class="flex {message.sender === user.user_id ? 'justify-end' : 'justify-start'}">
                <!-- Message Bubble -->
                <div class="max-w-[75%] p-3 rounded-lg {message.sender === user.user_id ? 'bg-[#25d366] text-white' : 'bg-[#E5E5E5] text-black'}">
                  <div class="flex items-center space-x-2">
                    {#if message.sender !== user.user_id}
                      <img
                        src="../../../user-icon.svg"
                        alt="User Avatar"
                        class="w-6 h-6 lg:w-8 lg:h-8 rounded-full"
                      />
                    {/if}
                    <div>
                      <p class="text-sm lg:text-base">{message.content}</p>
                      <span class="text-xs lg:text-sm {message.sender === user.user_id ? 'text-white/80' : 'text-black/60'}">{formatDate(message.created_at)}</span>

                      <!-- WhatsApp-style tick for read/unread message -->
                      {#if message.sender === user.user_id}
                      <div class="flex items-center space-x-1">
                        {#if message.isRead}
                          <!-- Display the white tick when message is read (for blue background) -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg>
                        {:else}
                          <!-- Display the light gray tick when message is unread -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg>
                        {/if}
                      </div>
                      {/if}
                    </div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
          <!-- Message Input -->
          <div class="mt-4 flex items-center space-x-2">
            <input
              type="text"
              class="flex-1 p-2 rounded-lg border border-gray-300"
              placeholder="Type a message..."
              bind:value={newMessage}
              on:keydown={(event: KeyboardEvent) => { if (event.key === 'Enter') sendMessage() }}
            />
            <button
              class="bg-[#25d366] text-white px-4 py-2 rounded-lg hover:bg-[#128C7E] whitespace-nowrap"
              on:click={sendMessage}
            >
              Send
            </button>
          </div>
        </div>
      {:else}
        <!-- No Chat Selected -->
        <div class="bg-white p-4 rounded-lg shadow h-full flex items-center justify-center">
          <h2 class="text-lg lg:text-xl font-bold text-[#25d366]">Select a chat to start messaging</h2>
        </div>
      {/if}
    </section>
    <footer class="mt-4 lg:mt-6 text-sm text-center text-gray-600 shrink-0">
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </main>
</div>
