import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface BottomCTAProps {
  onSignUp: () => void;
}

const BottomCTA: React.FC<BottomCTAProps> = ({ onSignUp }) => {
  return (
    <Card className="bg-gradient-to-r from-[#4A3726] to-[#634832] text-[#F9F6F4] border-none shadow-lg">
      <CardContent className="pt-8 pb-8 px-6">
        <h3 className="text-2xl font-bold mb-4 font-serif">
          Ready to join the community?
        </h3>
        <p className="mb-6 text-[#E8D7C9] text-lg">
          Sign up now to share your own coffee experiences, follow other
          enthusiasts, and discover new favorite spots.
        </p>
        <Button
          className="w-full bg-[#F9F6F4] text-[#4A3726] hover:bg-[#E8D7C9] transition-colors duration-300 text-lg py-3 rounded-full font-semibold"
          onClick={onSignUp}
        >
          Create Account
        </Button>
      </CardContent>
    </Card>
  );
};

export default BottomCTA;
