import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { UserAccount } from "@/api/models/UserAccount";
import { useSubscriptionMutations } from "../queries/useSubscriptionMutations";
import { useAuth } from "@/contexts/AuthContext";

interface AccountInfoProps {
  userAccount: UserAccount;
}

export function AccountInfo({ userAccount }: AccountInfoProps) {
  const { session } = useAuth();

  const subscriptionMutations = useSubscriptionMutations({
    onCheckoutSuccess: () => {
      alert("Redirecting to checkout...");
    },
    onCheckoutError: (error) => {
      alert(`Checkout failed: ${error.message}`);
    },
    onPortalSuccess: () => {
      alert("Redirecting to customer portal...");
    },
    onPortalError: (error) => {
      alert(`Portal failed: ${error.message}`);
    },
  });

  const handleUpgradeToPro = () => {
    subscriptionMutations.createCheckout.mutate({
      price_id: "price_1SKOuPLa4pEqShgojlivZTLc", // Your Stripe price ID
      success_url: `${window.location.origin}/dashboard?success=true`,
      cancel_url: `${window.location.origin}/dashboard?canceled=true`,
    });
  };

  const handleManageSubscription = () => {
    subscriptionMutations.createPortal.mutate({
      return_url: `${window.location.origin}/dashboard`,
    });
  };

  const accessToken = session?.access_token;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Account Information</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <p className="text-sm text-gray-600">Plan</p>
            <p className="font-medium">{userAccount.plan}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Credits</p>
            <p className="font-medium">{userAccount.credits}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Billing Customer ID</p>
            <p className="font-medium text-xs">
              {userAccount.billing_customer_id || "None"}
            </p>
          </div>
          {accessToken && (
            <div>
              <p className="text-sm text-gray-600">Access Token</p>
              <p className="font-medium text-xs break-all bg-gray-100 p-2 rounded-md">
                {accessToken}
              </p>
            </div>
          )}
        </div>

        <div className="mt-4 flex gap-2">
          {userAccount.plan === "free" ? (
            <Button onClick={handleUpgradeToPro}>
              {subscriptionMutations.createCheckout.isPending
                ? "Processing..."
                : "Upgrade to Pro"}
            </Button>
          ) : (
            <Button variant="outline" onClick={handleManageSubscription}>
              {subscriptionMutations.createPortal.isPending
                ? "Processing..."
                : "Manage Subscription"}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
