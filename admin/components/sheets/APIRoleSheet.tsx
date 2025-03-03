import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { CalendarIcon } from "lucide-react";
import { format } from "date-fns";
import { Input } from "@/components/ui/input";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Textarea } from "@/components/ui/textarea";
import { APIRoleType, CreateAPIRoleType, UpdateAPIRoleType } from "@/types/api_role";
import React, { useState } from "react";
import { IoCopyOutline, IoRefreshSharp, IoTrashOutline } from "react-icons/io5";
import { toast } from "sonner";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";

interface APIRoleSheetProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateAPIRoleType | UpdateAPIRoleType) => void;
  initialData?: APIRoleType | null;
  isSubmitting: boolean;
  onDelete?: (id: string) => Promise<void>;
  onRegenerateToken?: (id: string) => Promise<void>;
}

export default function APIRoleSheet({
  isOpen,
  onClose,
  onSubmit,
  initialData,
  isSubmitting,
  onDelete,
  onRegenerateToken,
}: APIRoleSheetProps) {
  const [name, setName] = useState(initialData?.Name || "");
  const [description, setDescription] = useState(initialData?.Description || "");
  const [expiresAt, setExpiresAt] = useState<Date | undefined>(
    initialData?.ExpiresAt ? new Date(initialData.ExpiresAt) : undefined
  );
  const [isRegeneratingKey, setIsRegeneratingKey] = useState(false);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [currentAPIKey, setCurrentAPIKey] = useState(initialData?.APIKey || "");

  // Update currentAPIKey when initialData changes
  useEffect(() => {
    setCurrentAPIKey(initialData?.APIKey || "");
  }, [initialData]);

  // Update expiresAt when initialData changes
  useEffect(() => {
    setExpiresAt(initialData?.ExpiresAt ? new Date(initialData.ExpiresAt) : undefined);
  }, [initialData]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // Only include fields that have values when updating
    if (initialData) {
        const updateData: UpdateAPIRoleType = {};
        if (name !== initialData.Name) updateData.name = name;
        if (description !== initialData.Description) updateData.description = description;
        // Always include expires_at in update
        updateData.expires_at = expiresAt?.toISOString() || null;
        onSubmit(updateData);
    } else {
        // For creation, require name and description
        onSubmit({
            name,
            description,
            expires_at: expiresAt?.toISOString() || null,
        });
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent className="bg-white dark:bg-black text-black dark:text-white">
        <SheetHeader>
          <SheetTitle>{initialData ? "\"" + initialData.Name + "\" token details" : "Create API Role"}</SheetTitle>
        </SheetHeader>
        {
          !initialData?.IsSystem && (
            <form onSubmit={handleSubmit} className="space-y-4 mt-4">
              <div className="space-y-2">
                <label
                  htmlFor="name"
                  className="text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                  Name
                </label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  defaultValue={initialData?.Name}
                  placeholder={initialData ? initialData.Name : "Enter role name"}
                  required={!initialData} // Only required when creating
                />
              </div>

              <div className="space-y-2">
                <label
                  htmlFor="description"
                  className="text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                  Description
                </label>
                <Textarea
                  id="description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  defaultValue={initialData?.Description}
                  placeholder={initialData ? initialData.Description : "Enter role description"}
                  required={!initialData} // Only required when creating
                />
              </div>

                    {/* Add expiration date field */}
                    <div className="space-y-2">
                      <label
                        htmlFor="expiresAt"
                        className="text-sm font-medium text-gray-700 dark:text-gray-300"
                      >
                        Expiration Date
                      </label>
                      <div>
                        {expiresAt ? (
                          <p className="text-sm text-gray-500">
                            This token will expire on {format(expiresAt, "PPP")}.
                          </p>
                        ):(
                          <p className="text-sm text-gray-500">
                            This token will not expire.
                          </p>
                        )}
                      </div>
                      <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant="outline"
                            className={`w-full justify-start text-left font-normal ${!expiresAt && "text-muted-foreground"}`}
                          >
                            <CalendarIcon className="mr-2 h-4 w-4" />
                            {expiresAt ? format(expiresAt, "PPP") : <span>Set expiration date</span>}
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                          <Calendar
                            mode="single"
                            selected={expiresAt}
                            onSelect={setExpiresAt}
                            initialFocus
                            disabled={(date) => date < new Date()}
                          />
                        </PopoverContent>
                      </Popover>
                      <div className="flex items-center justify-between mt-2">
                        <span className="text-sm text-gray-500">
                          {expiresAt ? "Token will expire on selected date" : "No expiration date set (permanent)"}
                        </span>
                        {expiresAt && (
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            className="text-red-500"
                            onClick={() => setExpiresAt(undefined)}
                          >
                            Remove expiration
                          </Button>
                        )}
                      </div>
                    </div>

                    <div className="flex justify-end space-x-2 pt-4">
                        <Button
                        type="button"
                        variant="outline"
                        onClick={onClose}
                        disabled={isSubmitting}
                        >
                        Cancel
                        </Button>
                        <Button type="submit" disabled={isSubmitting}>
                        {isSubmitting
                            ? initialData
                            ? "Saving..."
                            : "Creating..."
                            : initialData
                            ? "Save changes"
                            : "Create role"}
                        </Button>
                    </div>
                </form>
          )
        }
        {
            initialData?.Type != "public_user" && initialData && (
                <div className="flex flex-col justify-between pt-4">
                    <div className="text-sm font-medium text-gray-700 dark:text-gray-300">
                        Token
                    </div>
                    <div className="flex w-full my-2 gap-2 items-between">
                        <Button 
                            variant="secondary"
                            className="w-full"
                            onClick={() => {
                                navigator.clipboard.writeText(currentAPIKey);
                                toast.success("API key copied to clipboard");
                            }}
                        >
                            <IoCopyOutline className="w-3 h-3" />
                            Copy API key to clipboard
                        </Button>
                        <Button
                            variant="secondary"
                            onClick={async () => {
                                if (!initialData?.ID || !onRegenerateToken) return;
                                try {
                                    setIsRegeneratingKey(true);
                                    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles/${initialData.ID}`, {
                                        headers: {
                                            Authorization: `Bearer ${localStorage.getItem("token")}`,
                                        },
                                    });
                                    if (!response.ok) {
                                        throw new Error("Failed to fetch updated API key");
                                    }
                                    const data = await response.json();
                                    setCurrentAPIKey(data.APIKey);
                                    await onRegenerateToken(initialData.ID);
                                    toast.success("API key regenerated successfully");
                                // eslint-disable-next-line @typescript-eslint/no-unused-vars
                                } catch (error) {
                                    toast.error("Failed to regenerate API key");
                                } finally {
                                    setIsRegeneratingKey(false);
                                }
                            }}
                            disabled={isRegeneratingKey}
                        >
                            {isRegeneratingKey ? (
                                <Loader2 className="h-3 w-3 animate-spin" />
                            ) : (
                                <IoRefreshSharp className="w-3 h-3" />
                            )}
                        </Button>
                    </div>
                </div>
            )
        }
        {initialData && !initialData.IsSystem && (
            <div className="flex justify-end mt-4">
                <Button
                    variant="destructive"
                    onClick={() => setShowDeleteDialog(true)}
                    className="w-full"
                >
                    <IoTrashOutline className="w-3 h-3 mr-2" />
                    Delete API Role
                </Button>
            </div>
        )}

        <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Delete API Role</DialogTitle>
                    <DialogDescription>
                        Are you sure you want to delete this API role? This action cannot be undone.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <Button
                        variant="outline"
                        onClick={() => setShowDeleteDialog(false)}
                        disabled={isDeleting}
                    >
                        Cancel
                    </Button>
                    <Button
                        variant="destructive"
                        onClick={async () => {
                            if (!initialData?.ID || !onDelete) return;
                            try {
                                setIsDeleting(true);
                                await onDelete(initialData.ID);
                                toast.success("API role deleted successfully");
                                onClose();
                            // eslint-disable-next-line @typescript-eslint/no-unused-vars
                            } catch (error) {
                                toast.error("Failed to delete API role");
                            } finally {
                                setIsDeleting(false);
                                setShowDeleteDialog(false);
                            }
                        }}
                        disabled={isDeleting}
                    >
                        {isDeleting ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                Deleting...
                            </>
                        ) : (
                            "Delete"
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
      </SheetContent>
    </Sheet>
  );
}