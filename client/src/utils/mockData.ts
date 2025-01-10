import { Post } from "@/types";

export const mockFeedPosts: Post[] = [
  {
    id: 1,
    username: "coffeeexplorer",
    shopName: "File Systems of Coffee",
    location: "Downtown",
    imageUrl: "/api/placeholder/400/400",
    likes: 234,
    comments: 45,
    rating: 5,
    caption:
      "Perfect morning brew ☕️ The latte art here never disappoints! #coffeeart #morningcoffee",
    drinkName: "Hojicha Matcha Latte",
    drinkPrice: 5,
  },
  {
    id: 2,
    username: "beanconnoisseur",
    shopName: "Stereoscope Coffee",
    location: "Westside",
    imageUrl: "/api/placeholder/400/400",
    likes: 156,
    comments: 28,
    rating: 5,
    caption:
      "Their new single-origin Ethiopian beans are incredible! Notes of blueberry and dark chocolate.",
    drinkName: "Strawberry Matcha Latte w/ Sweetened Oat Milk",
    drinkPrice: 4,
  },
  {
    id: 3,
    username: "caffeinechaser",
    shopName: "Phin Smith Coffee",
    location: "Arts District",
    imageUrl: "/api/placeholder/400/400",
    likes: 312,
    comments: 67,
    rating: 4,
    caption:
      "Found my new favorite spot! The atmosphere is perfect for working and the coffee is exceptional.",
    drinkName: "Banana Coffee",
    drinkPrice: 3,
  },
];
