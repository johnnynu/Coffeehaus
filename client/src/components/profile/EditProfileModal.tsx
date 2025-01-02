import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { supabase } from "@/lib/supabase";
import { Check, X } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";

interface EditProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  profile: {
    username: string;
    display_name: string;
    bio: string | null;
  };
  onSave: (data: {
    username: string;
    display_name: string;
    bio: string;
  }) => Promise<void>;
}

const EditProfileModal: React.FC<EditProfileModalProps> = ({
  isOpen,
  onClose,
  profile,
  onSave,
}) => {
  const [formData, setFormData] = useState({
    username: "",
    display_name: "",
    bio: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [usernameStatus, setUsernameStatus] = useState<
    "available" | "taken" | "checking" | "initial"
  >("initial");
  const [debouncedUsername, setDebouncedUsername] = useState(formData.username);

  // Reset form data when modal opens
  useEffect(() => {
    if (isOpen && profile) {
      setFormData({
        username: profile.username,
        display_name: profile.display_name,
        bio: profile.bio || "",
      });
      setUsernameStatus("initial");
      setError(null);
    }
  }, [isOpen, profile]);

  // Debounce username changes
  useEffect(() => {
    const timer = setTimeout(() => {
      if (formData.username !== profile.username) {
        setDebouncedUsername(formData.username);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [formData.username, profile.username]);

  // Check username availability
  useEffect(() => {
    const checkUsername = async () => {
      if (debouncedUsername && debouncedUsername !== profile.username) {
        setUsernameStatus("checking");
        try {
          const { data } = await supabase
            .from("users")
            .select("username")
            .eq("username", debouncedUsername);

          setUsernameStatus(data && data.length > 0 ? "taken" : "available");
        } catch (error) {
          console.error("Error checking username:", error);
          setUsernameStatus("available");
        }
      } else {
        setUsernameStatus("initial");
      }
    };

    checkUsername();
  }, [debouncedUsername, profile.username]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      await onSave(formData);
      onClose();
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to update profile"
      );
      // Don't close modal on error so user can try again
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="space-y-4">
            <div>
              <label
                htmlFor="username"
                className="text-sm font-medium mb-1 block"
              >
                Username *
              </label>
              <div className="relative">
                <Input
                  id="username"
                  value={formData.username}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      username: e.target.value.trim().toLowerCase(),
                    }))
                  }
                  required
                  disabled={isSubmitting}
                  className={
                    usernameStatus === "available"
                      ? "pr-10 border-green-500 focus-visible:ring-green-500"
                      : usernameStatus === "taken"
                      ? "pr-10 border-red-500 focus-visible:ring-red-500"
                      : "pr-10"
                  }
                  placeholder="Enter a unique username"
                  pattern="^[a-z0-9_]+$"
                  title="Username can only contain lowercase letters, numbers, and underscores"
                />
                {usernameStatus === "checking" && (
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    <div className="animate-spin h-4 w-4 border-2 border-amber-800 rounded-full border-t-transparent" />
                  </div>
                )}
                {usernameStatus === "available" &&
                  formData.username !== profile.username && (
                    <Check className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-green-500" />
                  )}
                {usernameStatus === "taken" && (
                  <X className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-red-500" />
                )}
              </div>
              {usernameStatus === "taken" && (
                <p className="text-sm text-red-500 mt-1">
                  This username is already taken
                </p>
              )}
              <p className="text-xs text-gray-500 mt-1">
                Username can only contain lowercase letters, numbers, and
                underscores
              </p>
            </div>

            <div>
              <label
                htmlFor="displayName"
                className="text-sm font-medium mb-1 block"
              >
                Display Name *
              </label>
              <Input
                id="displayName"
                value={formData.display_name}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    display_name: e.target.value,
                  }))
                }
                required
                disabled={isSubmitting}
                placeholder="Your display name"
              />
            </div>

            <div>
              <label htmlFor="bio" className="text-sm font-medium mb-1 block">
                Bio
              </label>
              <Textarea
                id="bio"
                value={formData.bio}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, bio: e.target.value }))
                }
                className="h-24"
                disabled={isSubmitting}
                placeholder="Tell us about yourself..."
              />
            </div>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-md p-3">
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              className="bg-amber-800 text-white hover:bg-amber-700"
              disabled={
                isSubmitting ||
                usernameStatus === "taken" ||
                usernameStatus === "checking"
              }
            >
              {isSubmitting ? "Saving..." : "Save Changes"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default EditProfileModal;
