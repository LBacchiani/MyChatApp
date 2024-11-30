create table
  public."User" (
    user_id uuid not null default gen_random_uuid (),
    created_at timestamp with time zone not null default now(),
    username character varying not null,
    email character varying not null,
    constraint User_pkey primary key (user_id),
    constraint User_username_key unique (username),
    constraint User_email_key unique (email)
  ) tablespace pg_default;

create table
  public."Chat" (
    created_at timestamp with time zone not null default now(),
    participants uuid[] not null,
    owner_id uuid null,
    chat_id uuid not null default gen_random_uuid (),
    constraint Chat_pkey primary key (chat_id),
    constraint public_Chat_owner_id_fkey foreign key (owner_id) references "User" (user_id)
  ) tablespace pg_default;

create table
  public."Message" (
    message_id uuid not null default gen_random_uuid (),
    created_at timestamp with time zone not null default now(),
    content text not null,
    chat uuid not null,
    constraint Message_pkey primary key (message_id),
    constraint public_Message_chat_fkey foreign key (chat) references "Chat" (chat_id)
  ) tablespace pg_default;

create policy "Enable insert for users based on user_id" on "public"."User" as permissive for insert to public with check ((( SELECT auth.uid() AS uid) = user_id));
create policy "Users can read their info" on "public"."User" as permissive for select to public using ((( SELECT auth.uid() AS uid) = user_id));

create policy "Enable insert for authenticated users only" on "public"."Chat" as permissive for insert to authenticated with check (true);
create policy "Remove participants" on "public"."Chat" as permissive for update to authenticated using ((( SELECT auth.uid() AS uid) = owner_id));

create policy "Read Messages" on "public"."Message" as permissive for select to authenticated using ((EXISTS ( SELECT 1   FROM "Chat"  WHERE (("Chat".chat_id = "Message".chat) AND (auth.uid() = ANY ("Chat".participants))))));
create policy "Write Message" on "public"."Message" as permissive for insert to authenticated with check ((EXISTS ( SELECT 1   FROM "Chat"  WHERE (auth.uid() = ANY ("Chat".participants)))));