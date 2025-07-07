import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
} from '@mui/material';
import { GroupCreateRequest } from '../types/api';

interface CreateGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateGroup: (groupData: GroupCreateRequest) => Promise<void>;
  creatorId: string;
}

const CreateGroupDialog: React.FC<CreateGroupDialogProps> = ({
  open,
  onClose,
  onCreateGroup,
  creatorId,
}) => {
  const [newGroupData, setNewGroupData] = useState<Omit<GroupCreateRequest, 'image'>>({
    name: '',
    description: '',
    creator_id: '',
  });
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  const resetForm = () => {
    setNewGroupData({
      name: '',
      description: '',
      creator_id: '',
    });
    setSelectedImage(null);
    setImagePreview(null);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewGroupData({
      ...newGroupData,
      [e.target.name]: e.target.value,
    });
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // Validate file type
      const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
      if (!allowedTypes.includes(file.type)) {
        alert('Please select a valid image file (JPEG, PNG, GIF, or WebP)');
        return;
      }
      
      // Validate file size (10MB max)
      const maxSize = 10 * 1024 * 1024; // 10MB
      if (file.size > maxSize) {
        alert('File size must be less than 10MB');
        return;
      }
      
      setSelectedImage(file);
      
      // Create preview
      const reader = new FileReader();
      reader.onload = (e) => {
        setImagePreview(e.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleCreateGroup = async () => {
    if (!newGroupData.name.trim()) return;

    try {
      setIsCreating(true);
      
      // Create group data with optional image
      const groupDataWithCreator: GroupCreateRequest = {
        name: newGroupData.name,
        description: newGroupData.description,
        creator_id: creatorId,
        ...(selectedImage && { image: selectedImage }), // Only include image if selected
      };
      
      await onCreateGroup(groupDataWithCreator);
      resetForm();
      onClose();
    } catch (error) {
      console.error('Error creating group:', error);
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle>Create New Group</DialogTitle>
      <DialogContent>
        <TextField
          fullWidth
          name="name"
          label="Group Name"
          value={newGroupData.name}
          onChange={handleInputChange}
          required
          sx={{ mb: 2, mt: 1 }}
        />
        <TextField
          fullWidth
          name="description"
          label="Description"
          value={newGroupData.description}
          onChange={handleInputChange}
          multiline
          rows={3}
          sx={{ mb: 2 }}
        />
        <Box sx={{ mb: 2 }}>
          <Button
            variant="outlined"
            component="label"
            fullWidth
            sx={{ mb: 1 }}
          >
            {selectedImage ? 'Change Group Image' : 'Upload Group Image (Optional)'}
            <input
              type="file"
              hidden
              accept="image/*"
              onChange={handleImageChange}
            />
          </Button>
          {selectedImage && (
            <Typography variant="body2" color="text.secondary">
              Selected: {selectedImage.name}
            </Typography>
          )}
          {imagePreview && (
            <Box sx={{ mt: 1, textAlign: 'center' }}>
              <img
                src={imagePreview}
                alt="Preview"
                style={{
                  maxWidth: '200px',
                  maxHeight: '200px',
                  borderRadius: '4px',
                }}
              />
            </Box>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={isCreating}>
          Cancel
        </Button>
        <Button
          onClick={handleCreateGroup}
          variant="contained"
          color="secondary"
          disabled={!newGroupData.name.trim() || isCreating}
        >
          {isCreating ? 'Creating...' : 'Create Group'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateGroupDialog;
