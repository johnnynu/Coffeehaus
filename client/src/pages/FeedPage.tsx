import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { Coffee } from "lucide-react";
import { mockFeedPosts } from "@/utils/mockData";
import PostCard from "@/components/landing/PostCard";

interface UserProfile {
  username: string;
  display_name: string;
  profile_photo_id: string | null;
  photo_url?: string | null;
}

interface PhotoData {
  versions: {
    original: string;
  };
}

interface UserData extends Omit<UserProfile, "photo_url"> {
  photos: PhotoData;
}

const FeedPage: React.FC = () => {
  const navigate = useNavigate();
  const { user, session, signOut } = useAuth();
  const [profile, setProfile] = useState<UserProfile | null>(null);

  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!user?.id || !session?.access_token) return;

      try {
        const response = await fetch("http://localhost:8080/user", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${session.access_token}`,
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          console.error("Response status:", response.status);
          throw new Error("Failed to fetch user profile");
        }

        const userData = (await response.json()) as UserData;
        setProfile({
          ...userData,
          photo_url: userData.photos?.versions?.original || null,
        });
      } catch (error) {
        console.error("Error fetching user profile:", error);
      }
    };

    fetchUserProfile();
  }, [user?.id, session?.access_token]);

  const handleSignOut = async () => {
    try {
      await signOut();
      navigate("/");
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="fixed top-0 w-full bg-white border-b z-50">
        <div className="max-w-5xl mx-auto px-4 py-2 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <Coffee className="h-6 w-6 text-amber-800" />
            <h1 className="text-xl font-bold text-amber-800">Coffeehaus</h1>
          </div>
          <div className="flex items-center gap-4">
            {profile?.photo_url && (
              <img
                src={profile.photo_url}
                alt="Profile"
                className="w-8 h-8 rounded-full"
                crossOrigin="anonymous"
                referrerPolicy="no-referrer"
                onError={(e) => {
                  console.error("Error loading profile image:", e);
                  console.log("Attempted image URL:", profile.photo_url);
                }}
              />
            )}
            <Button
              variant="outline"
              className="border-amber-800 text-amber-800 hover:bg-amber-50"
              onClick={handleSignOut}
            >
              Sign Out
            </Button>
          </div>
        </div>
      </nav>

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
