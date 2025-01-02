import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import {
  Coffee,
  Grid3X3,
  Users,
  MapPin,
  Heart,
  MessageCircle,
} from "lucide-react";
import { mockFeedPosts } from "@/utils/mockData";
import EditProfileModal from "@/components/profile/EditProfileModal";
import useUserProfileStore from "@/store/userProfileStore";
import Navbar from "@/components/layout/Navbar";

const ProfilePage: React.FC = () => {
  const { username } = useParams<{ username: string }>();
  const navigate = useNavigate();
  const { session } = useAuth();
  const { profile, updateProfile, fetchProfile } = useUserProfileStore();
  const [isEditing, setIsEditing] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  // Fetch profile when username changes
  useEffect(() => {
    if (session?.access_token && username) {
      fetchProfile(session.access_token);
    }
  }, [username, session?.access_token, fetchProfile]);

  const handleProfileUpdate = async (data: {
    username: string;
    display_name: string;
    bio: string;
  }) => {
    if (!session?.access_token || !username) return;

    setIsUpdating(true);
    try {
      await updateProfile(session.access_token, username, data);
      setIsEditing(false);

      // If username changed, redirect to new profile URL
      if (data.username !== username) {
        navigate(`/profile/${data.username}`);
      } else {
        // Refetch profile data to ensure we have the latest
        await fetchProfile(session.access_token);
      }
    } catch (error) {
      console.error("Error updating profile:", error);
      throw error;
    } finally {
      setIsUpdating(false);
    }
  };

  const isOwnProfile = profile?.username === username;

  // Show loading state while profile is being updated
  if (isUpdating || !profile) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar isAuthenticated={true} />
        <div className="flex items-center justify-center h-[calc(100vh-64px)]">
          <div className="animate-spin h-8 w-8 border-4 border-amber-800 rounded-full border-t-transparent" />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#F9F6F4]">
      <Navbar isAuthenticated={true} />
      <main className="max-w-4xl mx-auto px-4 py-8 pt-20">
        {/* Profile Header */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
          <div className="flex flex-col md:flex-row items-start md:items-center gap-8 mb-8">
            {/* Profile Photo */}
            <div className="w-32 h-32 rounded-full overflow-hidden bg-[#E8D7C9] flex-shrink-0 border-4 border-[#D4B08C]">
              {profile.photo_url ? (
                <img
                  src={profile.photo_url}
                  alt={profile.display_name}
                  className="w-full h-full object-cover"
                  crossOrigin="anonymous"
                  referrerPolicy="no-referrer"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-4xl text-[#4A3726] font-serif">
                  {profile.display_name && profile.display_name.length > 0
                    ? profile.display_name[0].toUpperCase()
                    : "?"}
                </div>
              )}
            </div>

            {/* Profile Info */}
            <div className="flex-grow">
              <div className="flex items-center gap-4 mb-4">
                <h1 className="text-2xl font-semibold text-[#4A3726]">
                  {profile.username}
                </h1>
                {isOwnProfile && (
                  <Button
                    variant="outline"
                    className="border-[#4A3726] text-[#4A3726] hover:bg-[#4A3726] hover:text-[#F9F6F4]"
                    onClick={() => setIsEditing(true)}
                  >
                    Edit Profile
                  </Button>
                )}
              </div>

              <div className="flex gap-8 mb-4">
                <div className="text-sm text-[#634832]">
                  <span className="font-semibold">0</span> posts
                </div>
                <div className="text-sm text-[#634832]">
                  <span className="font-semibold">0</span> followers
                </div>
                <div className="text-sm text-[#634832]">
                  <span className="font-semibold">0</span> following
                </div>
              </div>

              <div>
                <h2 className="font-semibold mb-1 text-[#4A3726]">
                  {profile.display_name || profile.username}
                </h2>
                {profile.bio && (
                  <p className="text-sm whitespace-pre-wrap text-[#634832]">
                    {profile.bio}
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* Coffee Journey */}
          <div className="border-t border-[#E5E7EB] pt-4">
            <h3 className="text-lg font-semibold text-[#4A3726] mb-3">
              Coffee Journey
            </h3>
            <div className="flex justify-between">
              <div className="flex items-center gap-2 text-[#634832]">
                <Coffee className="w-5 h-5 text-[#967259]" />
                <span>12 coffee shops visited</span>
              </div>
              <div className="flex items-center gap-2 text-[#634832]">
                <MapPin className="w-5 h-5 text-[#967259]" />
                <span>5 cities explored</span>
              </div>
              <div className="flex items-center gap-2 text-[#634832]">
                <Users className="w-5 h-5 text-[#967259]" />
                <span>89 fellow coffee lovers</span>
              </div>
            </div>
          </div>
        </div>

        {/* Posts Grid */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="flex items-center justify-center gap-16 mb-6 text-sm font-medium">
            <button className="flex items-center gap-2 text-[#4A3726] border-b-2 border-[#4A3726] pb-2">
              <Grid3X3 className="w-4 h-4" />
              POSTS
            </button>
          </div>

          <div className="grid grid-cols-3 gap-4">
            {mockFeedPosts.map((post) => (
              <div
                key={post.id}
                className="aspect-square bg-[#E8D7C9] relative group cursor-pointer rounded-lg overflow-hidden"
              >
                <img
                  src={post.imageUrl}
                  alt={post.caption}
                  className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-[#4A3726] bg-opacity-0 group-hover:bg-opacity-30 transition-all duration-200 flex items-center justify-center">
                  <div className="text-[#F9F6F4] opacity-0 group-hover:opacity-100 transition-all duration-200 flex items-center gap-4">
                    <div className="flex items-center gap-1">
                      <Heart className="w-6 h-6" />
                      <span>{post.likes}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <MessageCircle className="w-6 h-6" />
                      <span>{post.comments}</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </main>

      {/* Edit Profile Modal */}
      {profile && (
        <EditProfileModal
          isOpen={isEditing}
          onClose={() => setIsEditing(false)}
          profile={profile}
          onSave={handleProfileUpdate}
        />
      )}
    </div>
  );
};

export default ProfilePage;
