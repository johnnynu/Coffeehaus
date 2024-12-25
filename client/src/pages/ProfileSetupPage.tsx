import React, { useState, ChangeEvent, FormEvent } from "react";
import { useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/contexts/AuthContext";

const ProfileSetupPage: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
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
      // TODO: Implement API call to save user profile
      console.log("Profile data:", formData);

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

        {user?.photoURL && (
          <div className="mb-6 flex justify-center">
            <img
              src={user.photoURL}
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
