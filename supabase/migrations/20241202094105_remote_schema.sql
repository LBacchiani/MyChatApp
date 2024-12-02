

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;


CREATE EXTENSION IF NOT EXISTS "pgsodium" WITH SCHEMA "pgsodium";






COMMENT ON SCHEMA "public" IS 'standard public schema';



CREATE EXTENSION IF NOT EXISTS "pg_graphql" WITH SCHEMA "graphql";






CREATE EXTENSION IF NOT EXISTS "pg_stat_statements" WITH SCHEMA "extensions";






CREATE EXTENSION IF NOT EXISTS "pgcrypto" WITH SCHEMA "extensions";






CREATE EXTENSION IF NOT EXISTS "pgjwt" WITH SCHEMA "extensions";






CREATE EXTENSION IF NOT EXISTS "supabase_vault" WITH SCHEMA "vault";






CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "extensions";





SET default_tablespace = '';

SET default_table_access_method = "heap";


CREATE TABLE IF NOT EXISTS "public"."Chat" (
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "participants" "uuid"[] NOT NULL,
    "owner_id" "uuid",
    "chat_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL
);


ALTER TABLE "public"."Chat" OWNER TO "postgres";


CREATE TABLE IF NOT EXISTS "public"."Message" (
    "message_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "content" "text" NOT NULL,
    "chat" "uuid" NOT NULL
);


ALTER TABLE "public"."Message" OWNER TO "postgres";


CREATE TABLE IF NOT EXISTS "public"."User" (
    "user_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "username" character varying NOT NULL,
    "email" character varying NOT NULL
);


ALTER TABLE "public"."User" OWNER TO "postgres";


ALTER TABLE ONLY "public"."Chat"
    ADD CONSTRAINT "chat_pkey" PRIMARY KEY ("chat_id");



ALTER TABLE ONLY "public"."Message"
    ADD CONSTRAINT "message_pkey" PRIMARY KEY ("message_id");



ALTER TABLE ONLY "public"."User"
    ADD CONSTRAINT "user_email_key" UNIQUE ("email");



ALTER TABLE ONLY "public"."User"
    ADD CONSTRAINT "user_pkey" PRIMARY KEY ("user_id");



ALTER TABLE ONLY "public"."User"
    ADD CONSTRAINT "user_username_key" UNIQUE ("username");



ALTER TABLE ONLY "public"."Chat"
    ADD CONSTRAINT "public_chat_owner_id_fkey" FOREIGN KEY ("owner_id") REFERENCES "public"."User"("user_id");



ALTER TABLE ONLY "public"."Message"
    ADD CONSTRAINT "public_message_chat_fkey" FOREIGN KEY ("chat") REFERENCES "public"."Chat"("chat_id");



ALTER TABLE "public"."Chat" ENABLE ROW LEVEL SECURITY;


CREATE POLICY "Enable insert for authenticated users only" ON "public"."Chat" FOR INSERT TO "authenticated" WITH CHECK (true);



CREATE POLICY "Insert new user" ON "public"."User" FOR INSERT WITH CHECK (true);



ALTER TABLE "public"."Message" ENABLE ROW LEVEL SECURITY;


CREATE POLICY "Read Messages" ON "public"."Message" FOR SELECT TO "authenticated" USING ((EXISTS ( SELECT 1
   FROM "public"."Chat"
  WHERE (("Chat"."chat_id" = "Message"."chat") AND ("auth"."uid"() = ANY ("Chat"."participants"))))));



CREATE POLICY "Remove participants" ON "public"."Chat" FOR UPDATE TO "authenticated" USING ((( SELECT "auth"."uid"() AS "uid") = "owner_id"));



CREATE POLICY "Update info" ON "public"."User" FOR UPDATE TO "authenticated" USING ((( SELECT "auth"."uid"() AS "uid") = "user_id"));



ALTER TABLE "public"."User" ENABLE ROW LEVEL SECURITY;


CREATE POLICY "Users can read their info" ON "public"."User" FOR SELECT TO "authenticated" USING ((( SELECT "auth"."uid"() AS "uid") = "user_id"));



CREATE POLICY "Write Message" ON "public"."Message" FOR INSERT TO "authenticated" WITH CHECK ((EXISTS ( SELECT 1
   FROM "public"."Chat"
  WHERE ("auth"."uid"() = ANY ("Chat"."participants")))));





ALTER PUBLICATION "supabase_realtime" OWNER TO "postgres";


GRANT USAGE ON SCHEMA "public" TO "postgres";
GRANT USAGE ON SCHEMA "public" TO "anon";
GRANT USAGE ON SCHEMA "public" TO "authenticated";
GRANT USAGE ON SCHEMA "public" TO "service_role";



































































































































































































GRANT ALL ON TABLE "public"."Chat" TO "anon";
GRANT ALL ON TABLE "public"."Chat" TO "authenticated";
GRANT ALL ON TABLE "public"."Chat" TO "service_role";



GRANT ALL ON TABLE "public"."Message" TO "anon";
GRANT ALL ON TABLE "public"."Message" TO "authenticated";
GRANT ALL ON TABLE "public"."Message" TO "service_role";



GRANT ALL ON TABLE "public"."User" TO "anon";
GRANT ALL ON TABLE "public"."User" TO "authenticated";
GRANT ALL ON TABLE "public"."User" TO "service_role";



ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES  TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES  TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES  TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES  TO "service_role";






ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS  TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS  TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS  TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS  TO "service_role";






ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES  TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES  TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES  TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES  TO "service_role";






























RESET ALL;
