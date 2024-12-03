create table
  public."User" (
    id uuid not null default gen_random_uuid (),
    created_at timestamp with time zone not null default now(),
    username character varying not null,
    email character varying not null,
    user_id uuid null,
    constraint user_pkey primary key (id),
    constraint user_email_key unique (email),
    constraint user_username_key unique (username),
    constraint User_user_id_fkey foreign key (user_id) references auth.users (id)
  ) tablespace pg_default;

create table
  public."Chat" (
    created_at timestamp with time zone not null default now(),
    participants uuid[] not null,
    owner_id uuid null,
    chat_id uuid not null default gen_random_uuid (),
    blocked boolean null default false,
    isGroup boolean null default false,
    constraint chat_pkey primary key (chat_id),
    constraint public_chat_owner_id_fkey foreign key (owner_id) references "User" (id)
  ) tablespace pg_default;

create table
  public."Message" (
    message_id uuid not null default gen_random_uuid (),
    created_at timestamp with time zone not null default now(),
    content text not null,
    chat uuid not null,
    sender uuid not null,
    constraint Message_pkey primary key (message_id),
    constraint public_Message_chat_fkey foreign key (chat) references "Chat" (chat_id),
    constraint Message_sender_fkey foreign key (sender) references "User" (user_id)
  ) tablespace pg_default;

create policy "Enable insert for authenticated users only" on "public"."User" as permissive for insert to public with check ( true );
create policy "Users can read their info" on "public"."User" as permissive for select to authenticated using ((( SELECT auth.uid() AS uid) = user_id));
create policy "Update info" on "public"."User" as permissive for update to authenticated using ((( SELECT auth.uid() AS uid) = user_id));

create policy "Select chat a user participates in" on "public"."Chat" as permissive for select to authenticated using (((SELECT auth.uid() AS uid) = ANY(participants)));
create policy "Enable insert for authenticated users only" on "public"."Chat" as permissive for insert to authenticated with check (true);
create policy "Modify participants" on "public"."Chat" as permissive for update to authenticated using ((isGroup = true and ( SELECT auth.uid() AS uid) = owner_id));
create policy "Block User" on "public"."Chat" as permissive for update to authenticated using (isGroup = false and (( SELECT auth.uid() AS uid) = owner_id));
create policy "Allow user to remove themselves from the chat" on "public"."Chat" as permissive for update to authenticated using (( SELECT auth.uid() AS uid) = ANY(participants) AND NOT (isGroup = true AND ( SELECT auth.uid() AS uid)  != owner_id) );


create policy "Read Messages" on "public"."Message" as permissive for select to authenticated using ((EXISTS ( SELECT 1   FROM "Chat"  WHERE (("Chat".chat_id = "Message".chat) AND (auth.uid() = ANY ("Chat".participants))))));
create policy "Write Message" on "public"."Message" as permissive for insert to authenticated with check ((EXISTS ( SELECT 1   FROM "Chat"  WHERE (auth.uid() = ANY ("Chat".participants)))));
create policy "Delete Message" on "public"."Message" as permissive for delete to authenticated using (( SELECT auth.uid() AS uid) = sender);
