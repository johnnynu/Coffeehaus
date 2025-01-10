import React from "react";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import {
  Heart,
  MessageCircle,
  MapPin,
  Coffee,
  Bookmark,
  MoreHorizontal,
} from "lucide-react";
import { Post } from "@/types";

interface PostCardProps {
  post: Post;
  onInteraction: () => void;
}

const PostCard: React.FC<PostCardProps> = ({ post, onInteraction }) => {
  const renderRating = (rating: number) => {
    const totalCoffees = 5;
    const fullCoffees = Math.round(rating);
    return (
      <div className="flex">
        {[...Array(totalCoffees)].map((_, index) => (
          <Coffee
            key={index}
            className={`w-3 h-3 ${
              index < fullCoffees
                ? "text-[#4A3726] fill-[#4A3726]"
                : "text-[#D4B08C]"
            }`}
          />
        ))}
      </div>
    );
  };

  return (
    <Card className="mb-4 bg-[#FAFAFA] border-[#E5E7EB] shadow-sm overflow-hidden">
      <CardContent className="p-0">
        {/* Post Header */}
        <div className="p-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-[#4A3726] to-[#967259] flex items-center justify-center text-[#F9F6F4] font-semibold text-xs">
              {post.username[0].toUpperCase()}
            </div>
            <div>
              <p className="font-semibold text-[#262626] text-sm">
                {post.username}
              </p>
              <div className="flex items-center gap-1 text-xs text-[#6B7280]">
                <MapPin className="w-3 h-3" />
                <span>{post.shopName}</span>
              </div>
            </div>
          </div>
          <button className="text-[#262626] hover:text-[#4A3726]">
            <MoreHorizontal className="w-5 h-5" />
          </button>
        </div>

        {/* Post Image */}
        <div className="aspect-square relative">
          <img
            src={post.imageUrl}
            alt="Coffee post"
            className="w-full h-full object-cover"
          />
        </div>

        {/* Post Actions */}
        <div className="p-3 bg-[#FAFAFA]">
          <div className="flex justify-between mb-2">
            <div className="flex gap-4">
              <button
                onClick={onInteraction}
                className="text-[#262626] hover:text-[#4A3726] transition-colors duration-200"
              >
                <Heart className="w-6 h-6" />
              </button>
              <button
                onClick={onInteraction}
                className="text-[#262626] hover:text-[#4A3726] transition-colors duration-200"
              >
                <MessageCircle className="w-6 h-6" />
              </button>
            </div>
            <button
              onClick={onInteraction}
              className="text-[#262626] hover:text-[#4A3726] transition-colors duration-200"
            >
              <Bookmark className="w-6 h-6" />
            </button>
          </div>

          {/* Post Details */}
          <div>
            <p className="text-sm font-semibold text-[#262626] mb-1">
              {post.likes} likes
            </p>
            <p className="text-sm text-[#262626]">
              <span className="font-semibold">{post.username}</span>{" "}
              {post.caption}
            </p>
            <div className="flex items-center gap-2 mt-1 text-xs text-[#6B7280]">
              {renderRating(post.rating)}
              <span className="font-medium">{post.drinkName}</span>
            </div>
            <p className="text-xs text-[#6B7280] mt-1">
              View all {post.comments} comments
            </p>
          </div>
        </div>
      </CardContent>

      <CardFooter className="bg-[#FAFAFA] border-t border-[#E5E7EB] p-3">
        <input
          type="text"
          placeholder="Add a comment..."
          className="w-full bg-transparent text-sm text-[#262626] placeholder-[#6B7280] focus:outline-none"
          onClick={onInteraction}
        />
      </CardFooter>
    </Card>
  );
};

export default PostCard;
