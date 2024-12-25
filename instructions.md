# Project Overview
You are building Coffeehaus, where users can discover and share experiences about coffee shops through an Instagram-style,
mobile-first web platform with location-based discovery, photo sharing, and social features. The app will be built using a modern tech stack combining:

Frontend: Vite.js and React with TypeScript, shadcn, tailwindcss, and an icon library of your choice
Backend: Go API with Chi
Authentication: Firebase Auth
Database: Supabase (PostgreSQL)
External Services:

Yelp API for coffee shop data
Google Maps for location features
Cloudinary for image processing
Google Cloud Platform for infrastructure

# Core Functionalities

## Priority 1: Core Authentication & Basic Data
### User Features
1. User Authentication (sign up with google, sign in with google, sign out) via Firebase Auth
    - When signing up/in with google, the user will be met with a google sign in popup to authorize the user
    - If the user is a first time user, they will be asked to enter a unique username, enter a display name, and enter an optional bio for themselves (we will use the profile photo from google as the avatar)
    - If the user is successful, the user should be able to login and see their feed
    - If the user fails authorization, there should be a alert/toast popup letting the user know that they failed to auth

2. Basic profile with username, avatar, and bio (optional)
    - When a new user signs up, they will be taken to a page to enter a new username, select their choice of avatar, and enter an optional bio for themselves
    - A user will have a basic profile page that is similar to how a user account from instagram works. Their profile will display their username, avatar, and bio (if they have one), as well as the posts they have created
    - Each post on the users profile should be similar to how instagram displays these posts, a square preview of the post that allows any auth'd user to click on which should lead to a modal of the post along with comments made on the said post
    - a username is unique to each user

### Coffee Shop Features
1. Integration with Yelp API for shop data
    - Shop data will be initially populated through Yelp API integration, providing verified business information including name, address, hours, and basic ratings
    - Users can discover coffee shops through a map interface, list view, or by searching with options to filter by distance, rating, or price level
    - Each coffee shop will have its own unique profile page displaying both Yelp data and Coffeehaus-specific content

## Priority 2: Basic Content & Viewing Features
### Social Features
1. Instagram-style photo posts with:
    - Single photo support initially (stored in Google Cloud Storage)
    - Location/shop tagging
    - Simple star rating (1-5)
    - Optional drink details
    - Basic captions

### Coffee Shop Features
1. Basic text search for shop names
    - Users can search for shops by name
    - Search results can be viewed in list format

2. Shop profiles with basic information
    - Each shop has a dedicated profile page featuring:
        * Basic information (hours, address, price level) from Yelp
        * Interactive map showing the shop's location
        * List of recent posts tagged at this location

## Priority 3: Social Engagement & Enhanced Discovery
### Social Features
1. Likes and comments
    - All auth users should be able to like and comment on a post. A user can only like a post once and a user can comment on a post as many times as they want. Users can also unlike a post

2. User following system
    - A user can follow another user once as well as unfollow a user. Once a user unfollows a user, they are allowed to follow that user again

### Coffee Shop Features
1. Location-based discovery
    - Users can see nearby coffee shops on an interactive map powered by Google Maps
    - Distance to each shop is displayed when browsing in list view
    - Users can filter shops within a specific radius of their current location
    - The app will request location permissions to provide personalized nearby recommendations

2. Dual Rating System (Yelp + Coffeehaus average user ratings)
    - Coffee shops will display both their Yelp rating and a Coffeehaus-specific rating
    - The Coffeehaus rating is calculated from user post ratings on Cofeehaus (1-5 stars) and updates in real-time as new ratings come in. This rating will provide how the coffee is to the user (whether they recommend it or not)
    - The Yelp rating comes from Yelp themselves as whenever we visit a coffee shop page on Yelp, we are given yelp's users ratings. This rating will provide the overall experience of the coffee shop
    - Users can see a breakdown of ratings and view the posts that contributed to the overall Coffeehaus score

## Priority 4: Advanced Features & Optimizations
### Social Features
1. Enhanced Instagram-style photo posts:
    - Multiple photo support up to 3
    - Captions with hashtag support

2. Activity feed
    - As of right now, a users activity feed will contain all posts made by anyone that uses the app. Whether they follow the user or not, similar to a user's personal feed

### User Features
1. Personal feed of users posts
    - As of right now, the users personal feed will contain all posts made by anyone that uses the app. Whether they follow the user or not

### Coffee Shop Features
1. Advanced search features
    - Autocomplete suggestions
    - Advanced filtering options including:
        * Operating hours (currently open)
        * Price level
        * Rating threshold
        * Distance range
        * Popular times
    - Map view integration for search results

2. Enhanced shop profiles
    - Additional features:
        * Photo gallery combining Yelp photos and user-submitted content
        * Quick statistics (total posts, average rating, popular drinks)
        * Option to save/bookmark the shop for later visits

# Doc

# Coffeehaus Data Models

## Core Entities

### User
The User model represents registered users of the Coffeehaus platform.
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    display_name TEXT,
    profile_photo_id UUID REFERENCES photos(id),
    bio TEXT,
    location POINT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Shop
The Shop model stores information about coffee shops, integrating data from Yelp with our platform's specific data.
```sql
CREATE TABLE shops (
    id UUID PRIMARY KEY,
    yelp_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    location POINT NOT NULL,
    address TEXT NOT NULL,
    yelp_rating DECIMAL(2,1),
    coffeehaus_rating DECIMAL(2,1),
    price_level TEXT,
    hours JSONB,
    last_sync TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Post
The Post model represents user-created content about their coffee experiences.
```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    shop_id UUID REFERENCES shops(id),
    caption TEXT,
    drink_name TEXT,
    drink_price DECIMAL,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Photo
The Photo model stores information about uploaded images, including their cloud storage locations and versions.
```sql
CREATE TABLE photos (
    id UUID PRIMARY KEY,
    post_id UUID REFERENCES posts(id),
    position INTEGER,
    versions JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Social Features

### Comment
The Comment model allows users to engage with posts through written responses.
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID REFERENCES posts(id),
    user_id UUID REFERENCES users(id),
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### PostLike (Junction Table)
Tracks which users have liked which posts, forming a many-to-many relationship.
```sql
CREATE TABLE post_likes (
    post_id UUID REFERENCES posts(id),
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);
```

### UserFollow (Junction Table)
Manages follow relationships between users, creating a self-referential many-to-many relationship.
```sql
CREATE TABLE user_follows (
    follower_id UUID REFERENCES users(id),
    following_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);
```

### SavedPost (Junction Table)
Allows users to bookmark posts for later reference.
```sql
CREATE TABLE saved_posts (
    user_id UUID REFERENCES users(id),
    post_id UUID REFERENCES posts(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);
```

### FavoriteShop (Junction Table)
Enables users to maintain a list of their favorite coffee shops.
```sql
CREATE TABLE favorite_shops (
    user_id UUID REFERENCES users(id),
    shop_id UUID REFERENCES shops(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, shop_id)
);
```

## Content Organization

### Hashtag
The Hashtag model supports content categorization and discovery.
```sql
CREATE TABLE hashtags (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
```

### PostHashtag (Junction Table)
Associates posts with hashtags, enabling content discovery and trending topics.
```sql
CREATE TABLE post_hashtags (
    post_id UUID REFERENCES posts(id),
    hashtag_id UUID REFERENCES hashtags(id),
    PRIMARY KEY (post_id, hashtag_id)
);
```

## Key Relationships

The data models form an interconnected system where:
- Users create posts about their experiences at coffee shops
- Posts contain photos and can be tagged with hashtags
- Users can follow other users, like posts, and favorite shops
- Comments and likes create engagement on posts
- Photos are associated with either posts or user profiles
- Shops maintain their own statistics and ratings


# Current File Structure
Coffeehaus
├── README.md
├── client
│   ├── README.md
│   ├── components.json
│   ├── eslint.config.js
│   ├── index.html
│   ├── package-lock.json
│   ├── package.json
│   ├── postcss.config.js
│   ├── public
│   │   └── vite.svg
│   ├── src
│   │   ├── App.css
│   │   ├── App.tsx
│   │   ├── assets
│   │   │   └── react.svg
│   │   ├── components
│   │   ├── hooks
│   │   ├── index.css
│   │   ├── lib
│   │   │   ├── supabase.ts
│   │   │   └── utils.ts
│   │   ├── main.tsx
│   │   ├── pages
│   │   ├── store
│   │   │   └── index.ts
│   │   ├── types
│   │   │   └── index.ts
│   │   ├── utils
│   │   └── vite-env.d.ts
│   ├── tailwind.config.js
│   ├── tsconfig.app.json
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   └── vite.config.ts
├── coffeehaus logo.png
├── instructions.md
└── server
    ├── cmd
    │   └── api
    └── internal
        ├── handlers
        ├── middleware
        ├── models
        └── utils

