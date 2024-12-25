export interface User {
  id: string;
  email: string;
  username: string | null;
  avatarUrl: string | null;
  bio: string | null;
}

export interface Post {
  id: number;
  username: string;
  shopName: string;
  location: string;
  imageUrl: string;
  likes: number;
  comments: number;
  rating: number;
  caption: string;
  drinkName: string;
  drinkPrice: number;
}

export interface Shop {
  id: string;
  yelpId: string;
  name: string;
  location: {
    lat: number;
    lng: number;
  };
  address: string;
  photos: string[];
  yelpRating: number;
  coffeehausRating: number | null;
  priceLevel: string;
  hours: {
    [key: string]: {
      open: string;
      close: string;
    };
  };
}
