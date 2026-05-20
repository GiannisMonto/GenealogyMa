-- Create cemeteries table
CREATE TABLE IF NOT EXISTS cemeteries (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100),
    address TEXT,
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    total_grave INTEGER DEFAULT 0,
    used_grave INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create graves table
CREATE TABLE IF NOT EXISTS graves (
    id BIGSERIAL PRIMARY KEY,
    cemetery_id BIGINT NOT NULL REFERENCES cemeteries(id) ON DELETE CASCADE,
    person_id BIGINT REFERENCES members(member_id) ON DELETE SET NULL,
    section VARCHAR(50) NOT NULL,
    row INTEGER NOT NULL,
    number INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved')),
    buried_name VARCHAR(200),
    buried_date DATE,
    buried_year INTEGER,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cemetery_id, section, row, number)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_cemeteries_province_city ON cemeteries(province, city);
CREATE INDEX IF NOT EXISTS idx_graves_cemetery_id ON graves(cemetery_id);
CREATE INDEX IF NOT EXISTS idx_graves_person_id ON graves(person_id);
CREATE INDEX IF NOT EXISTS idx_graves_status ON graves(status);