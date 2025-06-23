#!/bin/bash
set -e

# Check if seeding flag is enabled
if [ "$ENABLE_SEEDING" != "true" ]; then
  echo "Seeding is disabled. Set ENABLE_SEEDING=true to enable."
  exit 0
fi

echo "Starting data seeding..."

# Check if postgres client is installed
if ! [ -x "$(command -v psql)" ]; then
  echo "Error: psql is not installed." >&2
  exit 1
fi

# Check if database connection is available
until psql "$DB_DSN" -c '\l'; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - executing seed script"

# Check if users table exists and has data
user_count=$(psql "$DB_DSN" -tAc "SELECT COUNT(*) FROM users;")

if [ "$user_count" -gt "0" ]; then
  echo "Users table already has data. Skipping seeding to prevent duplicates."
  exit 0
fi

# Function to create a user and profile in one transaction
create_user_and_profile() {
    local user_id=$1
    local name=$2
    local username=$3
    local email=$4
    local user_type=$5
    local password_hash=$6 # Pre-hashed password for 'password'
    local firstname=$7
    local lastname=$8
    local title=$9
    local bio=${10}
    local faculty=${11}
    local program=${12}
    local degree=${13}
    local year=${14}
    local uni=${15}
    local skills=${16}
    
    # Create user
    psql "$DB_DSN" -c "
    INSERT INTO users (id, user_name, name, email, password_hash, activated, version, user_type, has_profile_created)
    VALUES ('$user_id', '$username', '$name', '$email', decode('$password_hash', 'hex'), true, 1, '$user_type', true);"

    # Create profile
    psql "$DB_DSN" -c "
    INSERT INTO user_profiles (user_id, firstname, lastname, title, bio, faculty, program, degree, year, uni)
    VALUES ('$user_id', '$firstname', '$lastname', '$title', '$bio', '$faculty', '$program', '$degree', '$year', '$uni');"

    # Add skills (from comma-separated list)
    IFS=',' read -ra ADDR <<< "$skills"
    for skill in "${ADDR[@]}"; do
        psql "$DB_DSN" -c "
        INSERT INTO user_skills (user_id, skill) 
        VALUES ('$user_id', '$(echo $skill | xargs)');"
    done
}

# Password hash for 'password123' - using SHA256 for demo (use bcrypt in production)
PASSWORD_HASH="243262243132243036784f6d36767a39436a6152766c7746354463627a4f345549546a386a37656241644e52707556373349385a636235644e3677"

# Create sample users with profiles
echo "Creating sample users and profiles..."

# Admin user
create_user_and_profile \
  "11d7b114-6ce3-49c9-a041-7acea988e16e" \
  "Admin User" \
  "admin" \
  "admin@example.com" \
  "admin" \
  "$PASSWORD_HASH" \
  "Admin" \
  "User" \
  "System Administrator" \
  "I manage the OpenConnect platform and help with technical issues." \
  "Information Technology" \
  "Computer Science" \
  "MSc" \
  "N/A" \
  "Open University of Sri Lanka" \
  "System Administration, Network Security, Cloud Computing, DevOps"

# Student user
create_user_and_profile \
  "22d7b114-6ce3-49c9-a041-7acea988e16e" \
  "John Doe" \
  "johndoe" \
  "john@example.com" \
  "student" \
  "$PASSWORD_HASH" \
  "John" \
  "Doe" \
  "Software Engineering Student" \
  "I'm a passionate software engineering student with interests in web development, AI, and blockchain technologies." \
  "Engineering and IT" \
  "Software Engineering" \
  "BSc" \
  "3" \
  "Open University of Sri Lanka" \
  "JavaScript, React, Go, Docker, PostgreSQL, AWS"

# Lecturer user
create_user_and_profile \
  "33d7b114-6ce3-49c9-a041-7acea988e16e" \
  "Sarah Johnson" \
  "sarahj" \
  "sarah@example.com" \
  "lecturer" \
  "$PASSWORD_HASH" \
  "Sarah" \
  "Johnson" \
  "Senior Lecturer in Computer Science" \
  "I specialize in artificial intelligence and machine learning, with a focus on natural language processing." \
  "Computing and Technology" \
  "Computer Science" \
  "PhD" \
  "N/A" \
  "Open University of Sri Lanka" \
  "Machine Learning, NLP, Python, TensorFlow, Research Methods, Data Analysis"

# Industry partner
create_user_and_profile \
  "44d7b114-6ce3-49c9-a041-7acea988e16e" \
  "Michael Chen" \
  "michaelc" \
  "michael@techcorp.com" \
  "industry" \
  "$PASSWORD_HASH" \
  "Michael" \
  "Chen" \
  "Senior Software Architect at TechCorp" \
  "I work with universities to provide real-world projects and mentoring opportunities for students." \
  "N/A" \
  "N/A" \
  "MSc" \
  "N/A" \
  "TechCorp Inc." \
  "System Architecture, Cloud Solutions, Microservices, Kubernetes, Mentoring"

# Create some ideas
echo "Creating sample ideas..."

psql "$DB_DSN" -c "
INSERT INTO ideas (id, title, description, user_id, category, tags, status, learning_outcome, recommended_level, created_at, updated_at, version)
VALUES 
  ('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 
   'AI-Powered Student Mentor System', 
   'A system that uses AI to provide personalized academic guidance to students, offering study suggestions, resources, and progress tracking.', 
   '22d7b114-6ce3-49c9-a041-7acea988e16e', 
   'Education Technology', 
   ARRAY['AI', 'Education', 'Mentorship', 'Student Support'], 
   'pending',
   'Students will receive timely guidance and support to improve academic performance.',
   'Intermediate',
   NOW(), NOW(), 1),
   
  ('b2c3d4e5-f6a7-8901-bcde-f23456789012', 
   'Blockchain-Based Academic Credential Verification', 
   'A blockchain solution that allows universities to issue tamper-proof digital certificates that can be instantly verified by employers.', 
   '33d7b114-6ce3-49c9-a041-7acea988e16e', 
   'Blockchain', 
   ARRAY['Blockchain', 'Education', 'Verification', 'Credentials'], 
   'approved',
   'Enhanced security and efficiency in academic credential verification.',
   'Advanced',
   NOW(), NOW(), 1),
   
  ('c3d4e5f6-a7b8-9012-cdef-345678901234', 
   'Machine Learning for Early Disease Detection', 
   'Research project using machine learning algorithms to analyze medical data for early detection of chronic diseases.', 
   '44d7b114-6ce3-49c9-a041-7acea988e16e', 
   'Health Research', 
   ARRAY['Machine Learning', 'Healthcare', 'Research', 'Disease Detection'], 
   'in_review',
   'Development of prediction models with practical applications in healthcare.',
   'Advanced',
   NOW(), NOW(), 1);"

echo "Data seeding completed successfully."
