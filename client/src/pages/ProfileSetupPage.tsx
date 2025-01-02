import React, { useState, ChangeEvent, FormEvent } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/contexts/AuthContext";
import { supabase } from "@/lib/supabase";
import { Coffee, FileText, User } from "lucide-react";

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
    <div className="min-h-screen bg-[#F9F6F4] flex items-center justify-center">
      <div className="max-w-md w-full p-8 bg-white rounded-lg shadow-md border border-[#D4B08C]">
        <div className="flex items-center justify-center mb-6">
          <Coffee className="w-10 h-10 text-[#4A3726]" />
        </div>
        <h1 className="text-3xl font-bold text-[#4A3726] mb-6 text-center">
          Brew Your Profile
        </h1>

        {avatarUrl && (
          <div className="mb-6 flex justify-center">
            <img
              src={avatarUrl}
              alt="Profile"
              className="w-24 h-24 rounded-full border-4 border-[#D4B08C]"
              crossOrigin="anonymous"
              referrerPolicy="no-referrer"
            />
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium mb-1 text-[#634832]"
            >
              Username *
            </label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#967259]" />
              <Input
                id="username"
                value={formData.username}
                onChange={handleInputChange}
                required
                placeholder="Choose a unique username"
                disabled={isSubmitting}
                className="pl-10 border-[#D4B08C] focus:ring-[#967259] text-[#4A3726]"
              />
            </div>
          </div>

          <div>
            <label
              htmlFor="displayName"
              className="block text-sm font-medium mb-1 text-[#634832]"
            >
              Display Name *
            </label>
            <div className="relative">
              <Coffee className="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#967259]" />
              <Input
                id="displayName"
                value={formData.displayName}
                onChange={handleInputChange}
                required
                placeholder="Your coffee connoisseur name"
                disabled={isSubmitting}
                className="pl-10 border-[#D4B08C] focus:ring-[#967259] text-[#4A3726]"
              />
            </div>
          </div>

          <div>
            <label
              htmlFor="bio"
              className="block text-sm font-medium mb-1 text-[#634832]"
            >
              Bio (Optional)
            </label>
            <div className="relative">
              <FileText className="absolute left-3 top-3 text-[#967259]" />
              <Textarea
                id="bio"
                value={formData.bio}
                onChange={handleInputChange}
                placeholder="Share your coffee journey..."
                className="pl-10 h-24 border-[#D4B08C] focus:ring-[#967259] text-[#4A3726]"
                disabled={isSubmitting}
              />
            </div>
          </div>

          <Button
            type="submit"
            className="w-full bg-[#4A3726] text-white hover:bg-[#634832] transition-colors duration-300"
            disabled={isSubmitting}
          >
            {isSubmitting
              ? "Brewing Your Profile..."
              : "Start Your Coffee Journey"}
          </Button>
        </form>
      </div>
    </div>
  );
};

export default ProfileSetupPage;
