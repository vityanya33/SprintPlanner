--
-- PostgreSQL database dump (с добавленным полем hours)
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';
SET default_table_access_method = heap;

--
-- TOC entry: tasks table with hours field
--

CREATE TABLE public.tasks (
                              id integer NOT NULL,
                              title text NOT NULL,
                              start_date date NOT NULL,
                              deadline date NOT NULL,
                              hours integer CHECK (hours > 0)
);

ALTER TABLE public.tasks OWNER TO postgres;

--
-- Остальная структура
--

CREATE TABLE public.task_users (
                                   task_id integer NOT NULL,
                                   user_id integer NOT NULL
);

ALTER TABLE public.task_users OWNER TO postgres;

CREATE TABLE public.users (
                              id integer NOT NULL,
                              name text NOT NULL,
                              role text
);

ALTER TABLE public.users OWNER TO postgres;

CREATE SEQUENCE public.tasks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.tasks_id_seq OWNER TO postgres;

ALTER TABLE public.tasks ALTER COLUMN id SET DEFAULT nextval('public.tasks_id_seq'::regclass);

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Данные
--

COPY public.tasks (id, title, start_date, deadline, hours) FROM stdin;
6	верстка	2025-07-16	2025-07-25	12
7	бэк	2025-07-25	2025-07-28	8
\.

COPY public.task_users (task_id, user_id) FROM stdin;
6	2
6	1
6	4
7	1
\.

COPY public.users (id, name, role) FROM stdin;
2	ывывы	ывыв
1	авваав11	вавав11
4	смсмсм смсс11	смсмсмсм
\.

--
-- Индексы и ограничения
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.task_users
    ADD CONSTRAINT task_users_pkey PRIMARY KEY (task_id, user_id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.task_users
    ADD CONSTRAINT task_users_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.task_users
    ADD CONSTRAINT task_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

--
-- Сброс счётчиков
--

SELECT pg_catalog.setval('public.tasks_id_seq', 7, true);
SELECT pg_catalog.setval('public.users_id_seq', 4, true);
