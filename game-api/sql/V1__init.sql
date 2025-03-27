CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- public.game_play definition

-- Drop table

-- DROP TABLE public.game_play;

CREATE TABLE public.game_play (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	status varchar NULL,
	"name" varchar NULL,
	current_frame int4 DEFAULT 1 NULL,
	current_roll int4 DEFAULT 1 NULL,
	current_user_index int4 DEFAULT 0 NULL,
	CONSTRAINT game_play_pk PRIMARY KEY (id)
);


-- public.player definition

-- Drop table

-- DROP TABLE public.player;

CREATE TABLE public.player (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	"name" varchar NOT NULL,
	CONSTRAINT player_pk PRIMARY KEY (id)
);


-- public.game_participant definition

-- Drop table

-- DROP TABLE public.game_participant;

CREATE TABLE public.game_participant (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	game_play_id uuid NULL,
	player_id uuid NULL,
	CONSTRAINT game_participant_pk PRIMARY KEY (id),
	CONSTRAINT game_participant_unique UNIQUE (player_id, game_play_id),
	CONSTRAINT game_participant_game_play_fk FOREIGN KEY (game_play_id) REFERENCES public.game_play(id),
	CONSTRAINT game_participant_player_fk FOREIGN KEY (player_id) REFERENCES public.player(id)
);


-- public.score definition

-- Drop table

-- DROP TABLE public.score;

CREATE TABLE public.score (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	game_play_id uuid NOT NULL,
	player_id uuid NOT NULL,
	score int4 NULL,
	frame int4 NULL,
	roll int4 NOT NULL,
	CONSTRAINT score_frame_unique UNIQUE (game_play_id, player_id, frame, roll),
	CONSTRAINT score_pk PRIMARY KEY (id),
	CONSTRAINT score_game_play_fk FOREIGN KEY (game_play_id) REFERENCES public.game_play(id),
	CONSTRAINT score_player_fk FOREIGN KEY (player_id) REFERENCES public.player(id)
);
