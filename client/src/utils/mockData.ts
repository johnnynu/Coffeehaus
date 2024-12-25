import { Post } from "@/types";

export const mockFeedPosts: Post[] = [
  {
    id: 1,
    username: "coffeeexplorer",
    shopName: "Artisan Coffee Co",
    location: "Downtown",
    imageUrl: "/api/placeholder/400/400",
    likes: 234,
    comments: 45,
    rating: 4.5,
    caption:
      "Perfect morning brew ☕️ The latte art here never disappoints! #coffeeart #morningcoffee",
    drinkName: "Oat Milk Latte",
    drinkPrice: 5,
  },
  {
    id: 2,
    username: "beanconnoisseur",
    shopName: "Roasters & Co",
    location: "Westside",
    imageUrl: "/api/placeholder/400/400",
    likes: 156,
    comments: 28,
    rating: 5,
    caption:
      "Their new single-origin Ethiopian beans are incredible! Notes of blueberry and dark chocolate.",
    drinkName: "Pour Over",
    drinkPrice: 4,
  },
  {
    id: 3,
    username: "caffeinechaser",
    shopName: "The Coffee Lab",
    location: "Arts District",
    imageUrl: "/api/placeholder/400/400",
    likes: 312,
    comments: 67,
    rating: 4.8,
    caption:
      "Found my new favorite spot! The atmosphere is perfect for working and the coffee is exceptional.",
    drinkName: "Cortado",
    drinkPrice: 3,
  },
];
