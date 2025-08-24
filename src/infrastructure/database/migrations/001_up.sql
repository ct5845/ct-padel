-- Players table
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Matches table
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    team1_player1_id INTEGER REFERENCES players(id),
    team1_player2_id INTEGER REFERENCES players(id),
    team2_player1_id INTEGER REFERENCES players(id),
    team2_player2_id INTEGER REFERENCES players(id),
    match_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sets table
CREATE TABLE sets (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(match_id, set_number)
);

-- Games table
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    set_id INTEGER REFERENCES sets(id) ON DELETE CASCADE,
    game_number INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(set_id, game_number)
);

-- Points table
CREATE TABLE points (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games(id) ON DELETE CASCADE,
    point_number INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, point_number)
);

-- Plays table (each shot/hit in a point)
CREATE TABLE plays (
    id SERIAL PRIMARY KEY,
    point_id INTEGER REFERENCES points(id) ON DELETE CASCADE,
    play_number INTEGER NOT NULL,
    player_id INTEGER REFERENCES players(id),
    ball_position_x INTEGER NOT NULL CHECK (ball_position_x >= 0 AND ball_position_x <= 6),
    ball_position_y INTEGER NOT NULL CHECK (ball_position_y >= 0 AND ball_position_y <= 11),
    result_type VARCHAR(50) CHECK (result_type IN ('unforced_error', 'error', 'no_return_winner')),
    hand_side VARCHAR(50) CHECK (hand_side IN ('forehand', 'backhand')),
    contact_type VARCHAR(50) CHECK (contact_type IN ('serve', 'groundstroke', 'volley', 'overhead')),
    shot_effect VARCHAR(50) CHECK (shot_effect IN ('flat', 'up', 'down', 'drop', 'smash')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(point_id, play_number)
);

-- Player positions for each play
CREATE TABLE player_positions (
    id SERIAL PRIMARY KEY,
    play_id INTEGER REFERENCES plays(id) ON DELETE CASCADE,
    player_id INTEGER REFERENCES players(id),
    position_x INTEGER NOT NULL CHECK (position_x >= 0 AND position_x <= 6),
    position_y INTEGER NOT NULL CHECK (position_y >= 0 AND position_y <= 11),
    UNIQUE(play_id, player_id)
);

-- Indexes for better performance
CREATE INDEX idx_sets_match_id ON sets(match_id);
CREATE INDEX idx_games_set_id ON games(set_id);
CREATE INDEX idx_points_game_id ON points(game_id);
CREATE INDEX idx_plays_point_id ON plays(point_id);
CREATE INDEX idx_play_positions_play_id ON player_positions(play_id);