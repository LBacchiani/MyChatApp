drop policy "Update info" on "public"."User";

alter table "public"."Chat" drop constraint "public_chat_owner_id_fkey";

alter table "public"."User" drop constraint "user_pkey";

drop index if exists "public"."user_pkey";

alter table "public"."Chat" add column "isGroup" boolean not null default false;

alter table "public"."User" add column "id" uuid not null default gen_random_uuid();

alter table "public"."User" alter column "user_id" drop default;

alter table "public"."User" alter column "user_id" drop not null;

CREATE UNIQUE INDEX user_pkey ON public."User" USING btree (id);

alter table "public"."User" add constraint "user_pkey" PRIMARY KEY using index "user_pkey";

alter table "public"."User" add constraint "User_user_id_fkey" FOREIGN KEY (user_id) REFERENCES auth.users(id) not valid;

alter table "public"."User" validate constraint "User_user_id_fkey";

alter table "public"."Chat" add constraint "public_chat_owner_id_fkey" FOREIGN KEY (owner_id) REFERENCES "User"(id) not valid;

alter table "public"."Chat" validate constraint "public_chat_owner_id_fkey";

create policy "Update info"
on "public"."User"
as permissive
for update
to authenticated
using ((( SELECT auth.uid() AS uid) = id));



