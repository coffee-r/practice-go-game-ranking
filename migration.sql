-- ユーザーテーブル
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

-- ゲームテーブル
CREATE TABLE games (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
)

-- ランキングテーブル
CREATE TABLE rankings (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    game_id INT NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_game_id FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
)

-- ユーザースコアテーブル
CREATE TABLE user_scores (
    ranking_id INT NOT NULL,
    user_id INT NOT NULL,
    score INT NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT pk_user_scores PRIMARY KEY (ranking_id, user_id)
    CONSTRAINT fk_ranking_id FOREIGN KEY (ranking_id) REFERENCES rankings(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)