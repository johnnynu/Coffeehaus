import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface BottomCTAProps {
  onSignUp: () => void;
}

const BottomCTA: React.FC<BottomCTAProps> = ({ onSignUp }) => {
  return (
    <Card className="bg-gradient-to-r from-amber-800 to-amber-700 text-white border-none">
      <CardContent className="pt-6">
        <h3 className="text-xl font-bold mb-2">Ready to join the community?</h3>
        <p className="mb-4 opacity-90">
          Sign up now to share your own coffee experiences, follow other
          enthusiasts, and discover new favorite spots.
        </p>
        <Button
          className="w-full bg-white text-amber-800 hover:bg-gray-100"
          onClick={onSignUp}
        >
          Create Account
        </Button>
      </CardContent>
    </Card>
  );
};

export default BottomCTA;
