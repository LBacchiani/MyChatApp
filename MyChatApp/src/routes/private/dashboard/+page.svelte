<script lang="ts">
  const SENDER_SERVICE = import.meta.env.VITE_SENDER_SERVICE;

  function formatDate(dateString: string) {
    const date = new Date(dateString);
    return  date.getDate() + "/" + date.getMonth() + "/" + date.getFullYear() + " | " + date.getHours() + ":" + date.getMinutes();
  }

  async function sendMessage() {
    if (!newMessage.trim()) return; // Prevent sending empty messages

    const messageData = {
      sender: user.user_id,
      receiver: selectedChat.user_id,
      chat_id: selectedChat.chat_id,
      content: newMessage,
    };

    try {
      const response = await fetch(SENDER_SERVICE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(messageData),
      });

      console.log(response)

      const result = await response.json();

      if (result.success) {
        selectedChat.messages = [...selectedChat.messages, { 
          sender: messageData.sender, 
          content: messageData.content, 
          created_at: new Date().toISOString() 
        }];
        newMessage = '';
      } else {
        alert('Message sending failed');
      }
    } catch (error) {
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
      const isDuplicate = chats.some(item => item.username === username2 && item.user_id === user_id);
      
      if (!isDuplicate) {
        const formData = new FormData();
        formData.append('participant1', user.user_id);
        formData.append('participant2', user_id);
        await fetch('?/createChat', {
          method: 'POST',
          body: formData,
        });   
        chats = [...chats, {username: username2, user_id: user_id, blocked: false, messages: []}];
      }
    }
  }

  export let data;
  let username2: string;
  let newMessage = '';
  const selectChat = (chat) => selectedChat = chat;
  $: ({user, chats} = data);
  $: selectedChat = chats[0];
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
        {#each chats as chat}
          <button 
            class="flex items-center p-2 w-full text-left bg-blue-500 rounded-lg hover:bg-blue-400 cursor-pointer"
            on:click={() => selectChat(chat)}
            aria-label="Select chat with {chat.username}"
            aria-pressed={selectedChat === chat ? 'true' : 'false'}
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
      <button
        class="w-full bg-gray-700 py-2 px-4 text-white rounded-lg hover:bg-gray-600 mb-2">
        Settings
      </button>
      <form method="post" action="?/logout">
        <button
          type="submit"
          class="w-full bg-red-600 py-2 px-4 text-white rounded-lg hover:bg-red-500">
          Logout
        </button>
      </form>
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

    <!-- Chat Conversation Section -->
    <section class="mt-6">
      {#if selectedChat}
        <!-- Conversation Section -->
        <div class="bg-white p-4 rounded-lg shadow-lg flex flex-col">
          <h2 class="text-xl font-bold text-blue-600 mb-4">{selectedChat.username}</h2>
          <div class="flex-grow overflow-y-auto space-y-4" style="max-height: 400px; overflow-y: scroll;">
            <!-- Display Messages -->
            {#each selectedChat.messages as message}
              <div class="flex {message.sender === user.user_id ? 'justify-end' : 'justify-start'}">
                <!-- Message Bubble -->
                <div class="max-w-xs p-3 rounded-lg {message.sender === user.user_id ? 'bg-blue-600 text-white' : 'bg-gray-200 text-black'}">
                  <div class="flex items-center space-x-2">
                    {#if message.sender !== user.user_id}
                      <img
                        src="../../../user-icon.svg"
                        alt="User Avatar"
                        class="w-8 h-8 rounded-full"
                      />
                    {/if}
                    <div>
                      <p class="text-sm">{message.content}</p>
                      <span class="text-xs {message.sender === user.user_id ? 'bg-blue-600 text-white' : 'bg-gray-200 text-black'}">{formatDate(message.created_at)}</span>
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
              class="flex-grow p-2 rounded-lg border border-gray-300"
              placeholder="Type a message..."
              bind:value={newMessage}
            />
            <button
              class="bg-blue-600 text-white p-2 rounded-lg hover:bg-blue-500"
              on:click={sendMessage}
            >
              Send
            </button>
          </div>
        </div>
      {:else}
        <!-- No Chat Selected -->
        <div class="bg-white p-4 rounded-lg shadow">
          <h2 class="text-xl font-bold text-blue-600">Select a chat to start messaging</h2>
        </div>
      {/if}
    </section>
    <footer>
      <p>&copy; {new Date().getFullYear()} MyChatApp. All Rights Reserved.</p>
    </footer>
  </main>
</div>
