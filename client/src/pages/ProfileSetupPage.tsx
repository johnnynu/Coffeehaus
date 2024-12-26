import React, { useState, ChangeEvent, FormEvent } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/contexts/AuthContext";
import { supabase } from "@/lib/supabase";

const ProfileSetupPage: React.FC = () => {
  const navigate = useNavigate();
  const { user, getAvatarUrl } = useAuth();
  const avatarUrl = getAvatarUrl();
  const [formData, setFormData] = useState({
    username: "",
    displayName: "",
    bio: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      if (!user) throw new Error("User not found");

      // check if username is available
      const { data: existingUser } = await supabase
        .from("users")
        .select("*")
        .eq("username", formData.username)
        .single();

      if (existingUser) {
        throw new Error("Username already taken");
      }

      // First create the photo entry if we have an avatar URL
      let photoId = null;
      if (avatarUrl) {
        const { data: photo, error: photoError } = await supabase
          .from("photos")
          .insert({
            versions: { original: avatarUrl },
            position: 0,
          })
          .select("id")
          .single();

        if (photoError) throw photoError;
        photoId = photo.id;
      }

      // create user profile
      const { error: insertError } = await supabase.from("users").insert({
        id: user.id,
        email: user.email!,
        username: formData.username,
        display_name: formData.displayName,
        bio: formData.bio || null,
        profile_photo_id: photoId,
      });

      if (insertError) {
        throw insertError;
      }

      // Navigate to feed page after successful submission
      navigate("/feed", { replace: true });
    } catch (error) {
      console.error("Error saving profile:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleInputChange = (
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="max-w-md w-full p-6 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-amber-800 mb-6">
          Complete Your Profile
        </h1>

        {avatarUrl && (
          <div className="mb-6 flex justify-center">
            <img
              src={avatarUrl}
              alt="Profile"
              className="w-24 h-24 rounded-full"
            />
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium mb-1"
            >
              Username *
            </label>
            <Input
              id="username"
              value={formData.username}
              onChange={handleInputChange}
              required
              placeholder="Choose a unique username"
              disabled={isSubmitting}
            />
          </div>

          <div>
            <label
              htmlFor="displayName"
              className="block text-sm font-medium mb-1"
            >
              Display Name *
            </label>
            <Input
              id="displayName"
              value={formData.displayName}
              onChange={handleInputChange}
              required
              placeholder="Your display name"
              disabled={isSubmitting}
            />
          </div>

          <div>
            <label htmlFor="bio" className="block text-sm font-medium mb-1">
              Bio (Optional)
            </label>
            <Textarea
              id="bio"
              value={formData.bio}
              onChange={handleInputChange}
              placeholder="Tell us about yourself..."
              className="h-24"
              disabled={isSubmitting}
            />
          </div>

          <Button
            type="submit"
            className="w-full bg-amber-800 text-white hover:bg-amber-700"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Setting up..." : "Complete Setup"}
          </Button>
        </form>
      </div>
    </div>
  );
};

export default ProfileSetupPage;
