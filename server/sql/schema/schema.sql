--
--


-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.6 (Homebrew)


--
-- Name: pgbouncer; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA pgbouncer;


--
-- Name: pg_graphql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_graphql WITH SCHEMA graphql;


--
-- Name: EXTENSION pg_graphql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_graphql IS 'pg_graphql: GraphQL support';


--
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA extensions;


--
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA extensions;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: supabase_vault; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS supabase_vault WITH SCHEMA vault;


--
-- Name: EXTENSION supabase_vault; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION supabase_vault IS 'Supabase Vault Extension';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA extensions;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: get_auth(text); Type: FUNCTION; Schema: pgbouncer; Owner: -
--

CREATE FUNCTION pgbouncer.get_auth(p_usename text) RETURNS TABLE(username text, password text)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $_$
begin
    raise debug 'PgBouncer auth request: %', p_usename;

    return query
    select 
        rolname::text, 
        case when rolvaliduntil < now() 
            then null 
            else rolpassword::text 
        end 
    from pg_authid 
    where rolname=$1 and rolcanlogin;
end;
$_$;


--
-- Name: handle_new_auth_user(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.handle_new_auth_user() RETURNS trigger
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
begin
  insert into public.user_accounts (id) values (new.id);
  return new;
end;
$$;


--
-- Name: tg_set_updated_at(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.tg_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  new.updated_at := now();
  return new;
end;
$$;




--
-- Name: ai_avatar_videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ai_avatar_videos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    filename text NOT NULL,
    thumbnail_filename text NOT NULL,
    duration integer,
    file_size bigint,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: credit_txns; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.credit_txns (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    request_id text NOT NULL,
    amount integer NOT NULL,
    status text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT credit_txns_amount_check CHECK ((amount > 0)),
    CONSTRAINT credit_txns_status_check CHECK ((status = ANY (ARRAY['reserved'::text, 'captured'::text, 'refunded'::text])))
);


--
-- Name: TABLE credit_txns; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.credit_txns IS 'Tracks credit transactions for idempotency and audit purposes';


--
-- Name: COLUMN credit_txns.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.id IS 'Unique transaction identifier';


--
-- Name: COLUMN credit_txns.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.user_id IS 'User who owns this transaction';


--
-- Name: COLUMN credit_txns.request_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.request_id IS 'Client-supplied idempotency key to prevent double-charging';


--
-- Name: COLUMN credit_txns.amount; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.amount IS 'Number of credits involved in this transaction (must be > 0)';


--
-- Name: COLUMN credit_txns.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.status IS 'Transaction status: reserved (pending), captured (completed), refunded (failed)';


--
-- Name: COLUMN credit_txns.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.created_at IS 'When the transaction was first created';


--
-- Name: COLUMN credit_txns.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.credit_txns.updated_at IS 'When the transaction was last updated';


--
-- Name: hooks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hooks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    generation_id uuid NOT NULL,
    prompt text NOT NULL,
    hook_text text NOT NULL,
    hook_index integer NOT NULL,
    credits_used integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT hooks_credits_used_check CHECK ((credits_used > 0)),
    CONSTRAINT hooks_hook_index_check CHECK ((hook_index >= 0))
);


--
-- Name: TABLE hooks; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.hooks IS 'Stores individual generated hooks for users';


--
-- Name: COLUMN hooks.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.id IS 'Unique hook identifier';


--
-- Name: COLUMN hooks.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.user_id IS 'User who generated this hook';


--
-- Name: COLUMN hooks.generation_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.generation_id IS 'Groups hooks from the same generation request';


--
-- Name: COLUMN hooks.prompt; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.prompt IS 'The prompt used to generate this hook';


--
-- Name: COLUMN hooks.hook_text; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.hook_text IS 'The actual hook text';


--
-- Name: COLUMN hooks.hook_index; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.hook_index IS 'Order of this hook within the generation (0-based)';


--
-- Name: COLUMN hooks.credits_used; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.credits_used IS 'Number of credits consumed for this generation';


--
-- Name: COLUMN hooks.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.created_at IS 'When the hook was generated';


--
-- Name: COLUMN hooks.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.hooks.updated_at IS 'When the record was last updated';


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL,
    applied_at timestamp with time zone DEFAULT now()
);


--
-- Name: user_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_accounts (
    id uuid NOT NULL,
    plan text DEFAULT 'free'::text NOT NULL,
    plan_started_at timestamp with time zone DEFAULT now() NOT NULL,
    plan_ends_at timestamp with time zone,
    billing_customer_id text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    credits integer DEFAULT 100 NOT NULL,
    CONSTRAINT user_accounts_credits_check CHECK ((credits >= 0)),
    CONSTRAINT user_accounts_plan_check CHECK ((plan = ANY (ARRAY['free'::text, 'pro'::text, 'enterprise'::text])))
);


--
-- Name: COLUMN user_accounts.credits; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_accounts.credits IS 'Number of credits available to the user (must be >= 0, default 100)';


--
-- Name: user_generated_videos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_generated_videos (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    ai_avatar_video_id uuid NOT NULL,
    overlay_text text NOT NULL,
    generated_video_filename character varying(255) NOT NULL,
    thumbnail_filename character varying(255) NOT NULL,
    status character varying(20) DEFAULT 'processing'::character varying,
    error_message text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: credit_txns credit_txns_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credit_txns
    ADD CONSTRAINT credit_txns_pkey PRIMARY KEY (id);


--
-- Name: credit_txns credit_txns_request_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credit_txns
    ADD CONSTRAINT credit_txns_request_id_key UNIQUE (request_id);


--
-- Name: hooks hooks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hooks
    ADD CONSTRAINT hooks_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_accounts user_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_accounts
    ADD CONSTRAINT user_accounts_pkey PRIMARY KEY (id);


--
-- Name: user_generated_videos user_generated_videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_generated_videos
    ADD CONSTRAINT user_generated_videos_pkey PRIMARY KEY (id);


--
-- Name: ai_avatar_videos videos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ai_avatar_videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY (id);


--
-- Name: idx_ai_avatar_videos_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ai_avatar_videos_created_at ON public.ai_avatar_videos USING btree (created_at DESC);


--
-- Name: idx_ai_avatar_videos_title; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ai_avatar_videos_title ON public.ai_avatar_videos USING btree (title);


--
-- Name: idx_credit_txns_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_credit_txns_created_at ON public.credit_txns USING btree (created_at);


--
-- Name: idx_credit_txns_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_credit_txns_status ON public.credit_txns USING btree (status);


--
-- Name: idx_credit_txns_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_credit_txns_user_id ON public.credit_txns USING btree (user_id);


--
-- Name: idx_hooks_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_hooks_created_at ON public.hooks USING btree (created_at);


--
-- Name: idx_hooks_generation_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_hooks_generation_id ON public.hooks USING btree (generation_id);


--
-- Name: idx_hooks_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_hooks_user_id ON public.hooks USING btree (user_id);


--
-- Name: idx_user_generated_videos_ai_avatar_video_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_generated_videos_ai_avatar_video_id ON public.user_generated_videos USING btree (ai_avatar_video_id);


--
-- Name: idx_user_generated_videos_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_generated_videos_created_at ON public.user_generated_videos USING btree (created_at);


--
-- Name: idx_user_generated_videos_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_generated_videos_status ON public.user_generated_videos USING btree (status);


--
-- Name: idx_user_generated_videos_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_generated_videos_user_id ON public.user_generated_videos USING btree (user_id);


--
-- Name: ai_avatar_videos set_updated_at_ai_avatar_videos; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at_ai_avatar_videos BEFORE UPDATE ON public.ai_avatar_videos FOR EACH ROW EXECUTE FUNCTION public.tg_set_updated_at();


--
-- Name: user_accounts set_updated_at_user_accounts; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at_user_accounts BEFORE UPDATE ON public.user_accounts FOR EACH ROW EXECUTE FUNCTION public.tg_set_updated_at();


--
-- Name: user_generated_videos set_updated_at_user_generated_videos; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_updated_at_user_generated_videos BEFORE UPDATE ON public.user_generated_videos FOR EACH ROW EXECUTE FUNCTION public.tg_set_updated_at();


--
-- Name: credit_txns credit_txns_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credit_txns
    ADD CONSTRAINT credit_txns_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_accounts(id) ON DELETE CASCADE;


--
-- Name: hooks hooks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hooks
    ADD CONSTRAINT hooks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_accounts(id) ON DELETE CASCADE;


--
-- Name: user_accounts user_accounts_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_accounts
    ADD CONSTRAINT user_accounts_id_fkey FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE;


--
-- Name: user_generated_videos user_generated_videos_ai_avatar_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_generated_videos
    ADD CONSTRAINT user_generated_videos_ai_avatar_video_id_fkey FOREIGN KEY (ai_avatar_video_id) REFERENCES public.ai_avatar_videos(id);


--
-- Name: user_generated_videos user_generated_videos_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_generated_videos
    ADD CONSTRAINT user_generated_videos_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_accounts(id);


--
-- Name: supabase_realtime; Type: PUBLICATION; Schema: -; Owner: -
--

CREATE PUBLICATION supabase_realtime WITH (publish = 'insert, update, delete, truncate');


--
-- Name: issue_graphql_placeholder; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER issue_graphql_placeholder ON sql_drop
         WHEN TAG IN ('DROP EXTENSION')
   EXECUTE FUNCTION extensions.set_graphql_placeholder();


--
-- Name: issue_pg_cron_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER issue_pg_cron_access ON ddl_command_end
         WHEN TAG IN ('CREATE EXTENSION')
   EXECUTE FUNCTION extensions.grant_pg_cron_access();


--
-- Name: issue_pg_graphql_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER issue_pg_graphql_access ON ddl_command_end
         WHEN TAG IN ('CREATE FUNCTION')
   EXECUTE FUNCTION extensions.grant_pg_graphql_access();


--
-- Name: issue_pg_net_access; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER issue_pg_net_access ON ddl_command_end
         WHEN TAG IN ('CREATE EXTENSION')
   EXECUTE FUNCTION extensions.grant_pg_net_access();


--
-- Name: pgrst_ddl_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER pgrst_ddl_watch ON ddl_command_end
   EXECUTE FUNCTION extensions.pgrst_ddl_watch();


--
-- Name: pgrst_drop_watch; Type: EVENT TRIGGER; Schema: -; Owner: -
--

CREATE EVENT TRIGGER pgrst_drop_watch ON sql_drop
   EXECUTE FUNCTION extensions.pgrst_drop_watch();


--
--


