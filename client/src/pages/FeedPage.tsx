import React, { useEffect } from "react";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";
import useUserProfileStore from "@/store/userProfileStore";
import { useAuth } from "@/contexts/AuthContext";
import Navbar from "@/components/layout/Navbar";

const FeedPage: React.FC = () => {
  const { session } = useAuth();
  const { profile, fetchProfile } = useUserProfileStore();

  useEffect(() => {
    if (session?.access_token) {
      fetchProfile(session.access_token);
    }
  }, [session?.access_token, fetchProfile]);

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar isAuthenticated={true} />

      <main className="pt-14">
        <div className="max-w-5xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-12 gap-8">
          <div className="md:col-span-4 md:sticky md:top-20 h-fit">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center gap-4 mb-4">
                {profile?.photo_url && (
                  <img
                    src={profile.photo_url}
                    alt="Profile"
                    className="w-16 h-16 rounded-full"
                    crossOrigin="anonymous"
                    referrerPolicy="no-referrer"
                    onError={(e) => {
                      console.error("Error loading profile image:", e);
                      console.log("Attempted image URL:", profile.photo_url);
                    }}
                  />
                )}
                <div>
                  <h2 className="font-semibold">{profile?.display_name}</h2>
                  <p className="text-gray-600 text-sm">@{profile?.username}</p>
                </div>
              </div>
              <p className="text-sm text-gray-600 mb-4">
                Welcome to your personalized coffee feed! Share your favorite
                coffee experiences and discover new spots.
              </p>
            </div>
          </div>

          <div className="md:col-span-8 space-y-6">
            {mockFeedPosts.map((post) => (
              <PostCard key={post.id} post={post} onInteraction={() => {}} />
            ))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default FeedPage;
