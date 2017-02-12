--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.1
-- Dumped by pg_dump version 9.6.1

-- Started on 2017-02-12 02:33:11

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12387)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2166 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 185 (class 1259 OID 16445)
-- Name: comments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE comments (
    id integer NOT NULL,
    created timestamp with time zone,
    author character varying,
    content character varying,
    post_id integer,
    parent_comment integer,
    edited timestamp with time zone NOT NULL,
    updated timestamp with time zone NOT NULL
);


ALTER TABLE comments OWNER TO postgres;

--
-- TOC entry 186 (class 1259 OID 16451)
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE comments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE comments_id_seq OWNER TO postgres;

--
-- TOC entry 2167 (class 0 OID 0)
-- Dependencies: 186
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE comments_id_seq OWNED BY comments.id;


--
-- TOC entry 187 (class 1259 OID 16453)
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE posts (
    id integer NOT NULL,
    author character varying NOT NULL,
    content character varying NOT NULL,
    created timestamp with time zone NOT NULL,
    edited timestamp with time zone NOT NULL,
    updated timestamp with time zone NOT NULL
);


ALTER TABLE posts OWNER TO postgres;

--
-- TOC entry 188 (class 1259 OID 16459)
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE posts_id_seq OWNER TO postgres;

--
-- TOC entry 2168 (class 0 OID 0)
-- Dependencies: 188
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE posts_id_seq OWNED BY posts.id;


--
-- TOC entry 193 (class 1259 OID 16489)
-- Name: securecookie; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE securecookie (
    id integer NOT NULL,
    hash bytea,
    block bytea
);


ALTER TABLE securecookie OWNER TO postgres;

--
-- TOC entry 192 (class 1259 OID 16487)
-- Name: securecookie_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE securecookie_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE securecookie_id_seq OWNER TO postgres;

--
-- TOC entry 2169 (class 0 OID 0)
-- Dependencies: 192
-- Name: securecookie_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE securecookie_id_seq OWNED BY securecookie.id;


--
-- TOC entry 189 (class 1259 OID 16461)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE sessions (
    key character varying NOT NULL,
    userid integer NOT NULL,
    logintime timestamp with time zone NOT NULL,
    lastseen timestamp with time zone
);


ALTER TABLE sessions OWNER TO postgres;

--
-- TOC entry 190 (class 1259 OID 16467)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    id integer NOT NULL,
    name character varying NOT NULL,
    email character varying,
    created timestamp with time zone NOT NULL,
    hash character varying NOT NULL,
    banned boolean
);


ALTER TABLE users OWNER TO postgres;

--
-- TOC entry 191 (class 1259 OID 16473)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO postgres;

--
-- TOC entry 2170 (class 0 OID 0)
-- Dependencies: 191
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- TOC entry 2028 (class 2604 OID 16475)
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY comments ALTER COLUMN id SET DEFAULT nextval('comments_id_seq'::regclass);


--
-- TOC entry 2029 (class 2604 OID 16476)
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY posts ALTER COLUMN id SET DEFAULT nextval('posts_id_seq'::regclass);


--
-- TOC entry 2031 (class 2604 OID 16492)
-- Name: securecookie id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY securecookie ALTER COLUMN id SET DEFAULT nextval('securecookie_id_seq'::regclass);


--
-- TOC entry 2030 (class 2604 OID 16477)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- TOC entry 2033 (class 2606 OID 16479)
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- TOC entry 2035 (class 2606 OID 16481)
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- TOC entry 2042 (class 2606 OID 16497)
-- Name: securecookie securecookie_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY securecookie
    ADD CONSTRAINT securecookie_pkey PRIMARY KEY (id);


--
-- TOC entry 2037 (class 2606 OID 16483)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (key);


--
-- TOC entry 2040 (class 2606 OID 16485)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (name, id);


--
-- TOC entry 2038 (class 1259 OID 16486)
-- Name: unique name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "unique name" ON users USING btree (name varchar_ops);

ALTER TABLE users CLUSTER ON "unique name";


-- Completed on 2017-02-12 02:33:12

--
-- PostgreSQL database dump complete
--

