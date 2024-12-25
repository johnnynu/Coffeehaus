import React from "react";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Heart, MessageCircle, Share2, MapPin, Coffee } from "lucide-react";
import { Post } from "@/types";

interface PostCardProps {
  post: Post;
  onInteraction: () => void;
}

const PostCard: React.FC<PostCardProps> = ({ post, onInteraction }) => {
  return (
    <Card className="mb-6">
      <CardContent className="p-0">
        {/* Post Header */}
        <div className="p-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-amber-100 flex items-center justify-center">
              {post.username[0].toUpperCase()}
            </div>
            <div>
              <p className="font-medium">{post.username}</p>
              <div className="flex items-center gap-1 text-sm text-gray-500">
                <MapPin className="w-3 h-3" />
                <span>{post.shopName}</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-1 text-amber-800">
            <Coffee className="w-4 h-4" />
            <span className="font-medium">{post.rating}</span>
          </div>
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
        <div className="p-4">
          <div className="flex gap-4 mb-3">
            <button
              onClick={onInteraction}
              className="text-gray-600 hover:text-amber-800"
            >
              <Heart className="w-6 h-6" />
            </button>
            <button
              onClick={onInteraction}
              className="text-gray-600 hover:text-amber-800"
            >
              <MessageCircle className="w-6 h-6" />
            </button>
            <button
              onClick={onInteraction}
              className="text-gray-600 hover:text-amber-800"
            >
              <Share2 className="w-6 h-6" />
            </button>
          </div>

          {/* Post Details */}
          <div>
            <p className="font-medium mb-1">{post.drinkName}</p>
            <p className="text-gray-600">{post.caption}</p>
            <p className="text-sm text-gray-500 mt-2">
              {post.likes} likes • {post.comments} comments
            </p>
          </div>
        </div>
      </CardContent>

      <CardFooter className="bg-gray-50 border-t">
        <button
          onClick={onInteraction}
          className="w-full text-gray-600 text-sm py-2"
        >
          Add a comment...
        </button>
      </CardFooter>
    </Card>
  );
};

export default PostCard;
