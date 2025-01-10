import { create } from "zustand";

interface UserProfile {
  username: string;
  display_name: string;
  bio: string | null;
  profile_photo_id: string | null;
  photo_url?: string | null;
}

interface UserProfileState {
  profile: UserProfile | null;
  isLoading: boolean;
  error: string | null;
  fetchProfile: (token: string) => Promise<void>;
  updateProfile: (
    token: string,
    username: string,
    data: Partial<UserProfile>
  ) => Promise<void>;
  reset: () => void;
}

const useUserProfileStore = create<UserProfileState>((set) => ({
  profile: null,
  isLoading: false,
  error: null,

  fetchProfile: async (token: string) => {
    console.log("Fetching profile...", new Date().toISOString());
    set({ isLoading: true, error: null });
    try {
      const response = await fetch("http://localhost:8080/user", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch user profile");
      }

      const userData = await response.json();
      set({
        profile: {
          ...userData,
          photo_url: userData.photos?.versions?.original || null,
        },
        isLoading: false,
      });
    } catch (error) {
      set({
        error:
          error instanceof Error ? error.message : "Failed to fetch profile",
        isLoading: false,
      });
    }
  },

  updateProfile: async (
    token: string,
    username: string,
    data: Partial<UserProfile>
  ) => {
    set({ isLoading: true, error: null });
    try {
      const response = await fetch(`http://localhost:8080/user/${username}`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        if (response.status === 409) {
          throw new Error("Username is already taken");
        }
        const errorData = await response.json().catch(() => null);
        throw new Error(
          errorData?.message || `Failed to update profile (${response.status})`
        );
      }

      const updatedProfile = await response.json();
      set({
        profile: {
          ...updatedProfile,
          photo_url: updatedProfile.photos?.versions?.original || null,
        },
        isLoading: false,
      });
    } catch (error) {
      set({
        error:
          error instanceof Error ? error.message : "Failed to update profile",
        isLoading: false,
      });
      throw error;
    }
  },

  reset: () => {
    set({ profile: null, isLoading: false, error: null });
  },
}));

export default useUserProfileStore;
