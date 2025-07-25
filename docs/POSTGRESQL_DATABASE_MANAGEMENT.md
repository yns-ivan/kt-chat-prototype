# PostgreSQL Database Management

## 🎯 **Recommended Tools**

### **1. pgAdmin (Most Popular for PostgreSQL)**
- **Type**: Web-based GUI
- **Download**: https://www.pgadmin.org/download/
- **Best for**: PostgreSQL-specific administration

### **2. DBeaver (Universal Database Tool)**
- **Type**: Desktop application
- **Download**: https://dbeaver.io/download/
- **Best for**: Multi-database support (PostgreSQL, MySQL, etc.)

### **3. TablePlus (Native Desktop)**
- **Type**: Native desktop application
- **Download**: https://tableplus.com/
- **Best for**: Clean, fast interface

## 🚀 **Quick Setup for Your Chat Project**

### **Connection Details**
```
Host: localhost
Port: 5432
Database: ktchat
Username: ktchat
Password: password
```

### **Option 1: pgAdmin Setup**

#### **1. Install pgAdmin**
```bash
# macOS (using Homebrew)
brew install --cask pgadmin4

# Or download from: https://www.pgadmin.org/download/
```

#### **2. Connect to Database**
1. Open pgAdmin
2. Right-click "Servers" → "Register" → "Server"
3. **General Tab**:
   - Name: `KT Chat Local`
4. **Connection Tab**:
   - Host: `localhost`
   - Port: `5432`
   - Database: `ktchat`
   - Username: `ktchat`
   - Password: `password`
5. Click "Save"

#### **3. Explore Your Database**
- **Tables**: `users`, `chat_rooms`, `messages`, `room_participants`
- **Views**: Any custom views
- **Functions**: Database functions

### **Option 2: DBeaver Setup**

#### **1. Install DBeaver**
```bash
# macOS (using Homebrew)
brew install --cask dbeaver-community

# Or download from: https://dbeaver.io/download/
```

#### **2. Connect to Database**
1. Open DBeaver
2. Click "New Database Connection"
3. Select "PostgreSQL"
4. **Connection Settings**:
   - Host: `localhost`
   - Port: `5432`
   - Database: `ktchat`
   - Username: `ktchat`
   - Password: `password`
5. Click "Test Connection" → "Finish"

### **Option 3: TablePlus Setup**

#### **1. Install TablePlus**
```bash
# macOS (using Homebrew)
brew install --cask tableplus

# Or download from: https://tableplus.com/
```

#### **2. Connect to Database**
1. Open TablePlus
2. Click "Create a new connection"
3. Select "PostgreSQL"
4. **Connection Details**:
   - Name: `KT Chat Local`
   - Host: `localhost`
   - Port: `5432`
   - Database: `ktchat`
   - User: `ktchat`
   - Password: `password`
5. Click "Connect"

## 🛠️ **Common Operations**

### **View Tables**
```sql
-- List all tables
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public';

-- View table structure
\d table_name
```

### **Query Messages**
```sql
-- Get recent messages with user info
SELECT 
    m.id,
    m.content,
    u.username,
    m.created_at
FROM messages m
JOIN users u ON m.user_id = u.id
ORDER BY m.created_at DESC
LIMIT 10;
```

### **Query Users**
```sql
-- Get all users
SELECT id, username, email, created_at
FROM users
ORDER BY created_at DESC;
```

### **Query Chat Rooms**
```sql
-- Get all chat rooms with participant count
SELECT 
    cr.id,
    cr.name,
    cr.description,
    COUNT(rp.id) as participant_count
FROM chat_rooms cr
LEFT JOIN room_participants rp ON cr.id = rp.room_id AND rp.left_at IS NULL
GROUP BY cr.id, cr.name, cr.description;
```

## 🔧 **Docker Database Access**

### **Direct Container Access**
```bash
# Connect to PostgreSQL container
docker exec -it ktchat-postgres psql -U ktchat -d ktchat

# List tables
\dt

# View table structure
\d users
\d messages
\d chat_rooms
```

### **Backup Database**
```bash
# Create backup
docker exec ktchat-postgres pg_dump -U ktchat ktchat > backup.sql

# Restore backup
docker exec -i ktchat-postgres psql -U ktchat ktchat < backup.sql
```

## 📊 **Useful Queries for Development**

### **Check Message Counts**
```sql
-- Messages per room
SELECT 
    cr.name as room_name,
    COUNT(m.id) as message_count
FROM chat_rooms cr
LEFT JOIN messages m ON cr.id = m.room_id
GROUP BY cr.id, cr.name
ORDER BY message_count DESC;
```

### **Check User Activity**
```sql
-- Users with most messages
SELECT 
    u.username,
    COUNT(m.id) as message_count,
    MAX(m.created_at) as last_message
FROM users u
LEFT JOIN messages m ON u.id = m.user_id
GROUP BY u.id, u.username
ORDER BY message_count DESC;
```

### **Check Room Participation**
```sql
-- Active participants per room
SELECT 
    cr.name as room_name,
    COUNT(rp.id) as active_participants
FROM chat_rooms cr
LEFT JOIN room_participants rp ON cr.id = rp.room_id AND rp.left_at IS NULL
GROUP BY cr.id, cr.name
ORDER BY active_participants DESC;
```

## 🎯 **Recommendation**

For your PostgreSQL setup, I recommend **pgAdmin** because:
- ✅ **PostgreSQL-specific**: Optimized for PostgreSQL features
- ✅ **Familiar workflow**: Similar to MySQL Workbench
- ✅ **Free and open-source**: No licensing costs
- ✅ **Full-featured**: Complete database administration
- ✅ **Web-based**: Accessible from any browser

Start with pgAdmin, and if you need a more universal tool later, DBeaver is an excellent alternative! 