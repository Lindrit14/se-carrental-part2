import { PageHeader } from '@/components/layout/PageHeader';
import { ProfileCard } from '@/features/profile/ProfileCard';
import { DeleteAccountDialog } from '@/features/profile/DeleteAccountDialog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';

export function ProfilePage() {
  return (
    <div className="flex flex-col gap-8">
      <PageHeader title="Profile" description="Manage your account." />
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-[2fr_1fr]">
        <ProfileCard />
        <Card>
          <CardHeader>
            <CardTitle>Danger zone</CardTitle>
            <CardDescription>
              Deleting your account is permanent. Existing bookings remain in the booking
              service but are anonymized.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <DeleteAccountDialog />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
