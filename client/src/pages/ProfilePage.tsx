import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { Grid3X3 } from "lucide-react";
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
    <div className="min-h-screen bg-gray-50">
      <Navbar isAuthenticated={true} />
      <main className="max-w-4xl mx-auto px-4 py-8 pt-14">
        {/* Profile Header */}
        <div className="flex flex-col md:flex-row items-start md:items-center gap-8 mb-8">
          {/* Profile Photo */}
          <div className="w-32 h-32 rounded-full overflow-hidden bg-gray-200 flex-shrink-0">
            {profile.photo_url ? (
              <img
                src={profile.photo_url}
                alt={profile.display_name}
                className="w-full h-full object-cover"
                crossOrigin="anonymous"
                referrerPolicy="no-referrer"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-4xl text-gray-400">
                {profile.display_name && profile.display_name.length > 0
                  ? profile.display_name[0].toUpperCase()
                  : "?"}
              </div>
            )}
          </div>

          {/* Profile Info */}
          <div className="flex-grow">
            <div className="flex items-center gap-4 mb-4">
              <h1 className="text-xl font-semibold">{profile.username}</h1>
              {isOwnProfile && (
                <Button
                  variant="outline"
                  className="border-gray-300"
                  onClick={() => setIsEditing(true)}
                >
                  Edit Profile
                </Button>
              )}
            </div>

            <div className="flex gap-8 mb-4">
              <div className="text-sm">
                <span className="font-semibold">0</span> posts
              </div>
              <div className="text-sm">
                <span className="font-semibold">0</span> followers
              </div>
              <div className="text-sm">
                <span className="font-semibold">0</span> following
              </div>
            </div>

            <div>
              <h2 className="font-semibold mb-1">
                {profile.display_name || profile.username}
              </h2>
              {profile.bio && (
                <p className="text-sm whitespace-pre-wrap">{profile.bio}</p>
              )}
            </div>
          </div>
        </div>

        {/* Posts Grid */}
        <div className="border-t pt-4">
          <div className="flex items-center justify-center gap-16 mb-4 text-sm font-medium">
            <button className="flex items-center gap-2 text-amber-800 border-t-2 border-amber-800 pt-4 -mt-4">
              <Grid3X3 className="w-4 h-4" />
              POSTS
            </button>
          </div>

          <div className="grid grid-cols-3 gap-1 md:gap-6">
            {mockFeedPosts.map((post) => (
              <div
                key={post.id}
                className="aspect-square bg-gray-200 relative group cursor-pointer"
              >
                <img
                  src={post.imageUrl}
                  alt={post.caption}
                  className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-200" />
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
