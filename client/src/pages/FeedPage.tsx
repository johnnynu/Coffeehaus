import React, { useCallback, useEffect, useState } from "react";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";
import useUserProfileStore from "@/store/userProfileStore";
import { useAuth } from "@/contexts/AuthContext";
import Navbar from "@/components/layout/Navbar";
import { Button } from "@/components/ui/button";
import { Coffee, MapPin, PlusCircle, Users } from "lucide-react";

const FeedPage: React.FC = () => {
  const { session } = useAuth();
  const { profile, fetchProfile } = useUserProfileStore();
  const [activeTab, setActiveTab] = useState<"following" | "discover">(
    "discover"
  );

  const memoizedFetchProfile = useCallback(() => {
    if (session?.access_token) {
      console.log("Triggering fetch...", new Date().toISOString());
      fetchProfile(session.access_token);
    }
  }, [session?.access_token, fetchProfile]);

  useEffect(() => {
    memoizedFetchProfile();
  }, [memoizedFetchProfile]);

  return (
    <div className="min-h-screen bg-[#F9F6F4]">
      <Navbar isAuthenticated={true} />

      <main className="pt-16">
        <div className="max-w-6xl mx-auto px-4 py-6 grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Sidebar */}
          <div className="lg:col-span-3 lg:sticky lg:top-24 h-fit space-y-6">
            <div className="bg-white rounded-lg shadow-sm p-6">
              <div className="flex items-center gap-4 mb-4">
                {profile?.photo_url ? (
                  <img
                    src={profile.photo_url}
                    alt="Profile"
                    className="w-16 h-16 rounded-full border-2 border-[#D4B08C]"
                    crossOrigin="anonymous"
                    referrerPolicy="no-referrer"
                    onError={(e) => {
                      console.error("Error loading profile image:", e);
                      console.log("Attempted image URL:", profile.photo_url);
                    }}
                  />
                ) : (
                  <div className="w-16 h-16 rounded-full bg-[#4A3726] flex items-center justify-center text-[#F9F6F4] text-2xl font-bold">
                    {profile?.display_name?.[0].toUpperCase()}
                  </div>
                )}
                <div>
                  <h2 className="font-semibold text-[#4A3726]">
                    {profile?.display_name}
                  </h2>
                  <p className="text-[#6B7280] text-sm">@{profile?.username}</p>
                </div>
              </div>
              <p className="text-sm text-[#634832] mb-4">
                Welcome to your personalized coffee feed! Share your favorite
                coffee experiences and discover new spots.
              </p>
              <div className="flex flex-col gap-2">
                <Button
                  className="bg-[#4A3726] text-[#F9F6F4] hover:bg-[#634832] py-2"
                  onClick={() => {}}
                >
                  <PlusCircle className="w-4 h-4 mr-2 shrink-0" />
                  New Post
                </Button>
                <Button
                  className="bg-[#4A3726] text-[#F9F6F4] hover:bg-[#634832] py-2"
                  onClick={() => {}}
                >
                  <MapPin className="w-4 h-4 mr-2 shrink-0" />
                  Discover Nearby
                </Button>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-sm p-6">
              <h3 className="font-semibold text-[#4A3726] mb-4">
                Your Coffee Journey
              </h3>
              <ul className="space-y-3">
                <li className="flex items-center text-[#634832]">
                  <Coffee className="w-5 h-5 mr-2 text-[#967259]" />
                  <span>12 coffee shops visited</span>
                </li>
                <li className="flex items-center text-[#634832]">
                  <MapPin className="w-5 h-5 mr-2 text-[#967259]" />
                  <span>5 cities explored</span>
                </li>
                <li className="flex items-center text-[#634832]">
                  <Users className="w-5 h-5 mr-2 text-[#967259]" />
                  <span>89 fellow coffee lovers</span>
                </li>
              </ul>
            </div>
          </div>

          {/* Main Feed */}
          <div className="lg:col-span-9 space-y-6">
            <div className="bg-white rounded-lg shadow-sm p-4">
              <div className="flex border-b border-[#E5E7EB]">
                <button
                  className={`flex-1 py-2 text-center font-medium ${
                    activeTab === "discover"
                      ? "text-[#4A3726] border-b-2 border-[#4A3726]"
                      : "text-[#6B7280]"
                  }`}
                  onClick={() => setActiveTab("discover")}
                >
                  Discover
                </button>
                <button
                  className={`flex-1 py-2 text-center font-medium ${
                    activeTab === "following"
                      ? "text-[#4A3726] border-b-2 border-[#4A3726]"
                      : "text-[#6B7280]"
                  }`}
                  onClick={() => setActiveTab("following")}
                >
                  Following
                </button>
              </div>
            </div>
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
