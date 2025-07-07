import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import { ActivityCreateRequest, Group } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface AddActivityDialogProps {
  open: boolean;
  onClose: () => void;
  onActivityAdded: () => void;
}

const AddActivityDialog: React.FC<AddActivityDialogProps> = ({
  open,
  onClose,
  onActivityAdded,
}) => {
  const { user } = useAuth();
  const [formData, setFormData] = useState<Omit<ActivityCreateRequest, 'image'>>({
    title: '',
    description: '',
    date: new Date().toISOString().split('T')[0], // Today's date
    group_id: '',
    creator_id: user?.id || '', // Include creator_id from auth context
  });
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open && user) {
      // Update creator_id when user changes or dialog opens
      setFormData(prevData => ({
        ...prevData,
        creator_id: user.id
      }));
      loadUserGroups();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user]);

  const loadUserGroups = async () => {
    if (!user) return;
    
    try {
      const userGroups = await apiService.getUserGroups(user.id);
      setGroups(userGroups);
    } catch (error) {
      console.error('Error loading user groups:', error);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const name = e.target.name as keyof Omit<ActivityCreateRequest, 'image'>;
    const value = e.target.value as string;
    
    setFormData({
      ...formData,
      [name]: value,
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

  const handleSelectChange = (e: any) => {
    const name = e.target.name as keyof Omit<ActivityCreateRequest, 'image'>;
    const value = e.target.value as string;
    
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !selectedImage) return;

    try {
      setLoading(true);
      
      const activityData: ActivityCreateRequest = {
        ...formData,
        image: selectedImage,
      };
      
      await apiService.createActivity(activityData);
      
      // Reset form
      setFormData({
        title: '',
        description: '',
        date: new Date().toISOString().split('T')[0],
        group_id: '',
        creator_id: user.id, // Reset creator_id to current user
      });
      setSelectedImage(null);
      setImagePreview(null);
      
      onActivityAdded();
      onClose();
    } catch (error) {
      console.error('Error creating activity:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="md"
      fullWidth
    >
      <DialogTitle>
        <Typography variant="h2">
          Log Your Activity
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Share what you've accomplished today
        </Typography>
      </DialogTitle>
      
      <DialogContent>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            fullWidth
            name="title"
            label="Activity Title"
            value={formData.title}
            onChange={handleInputChange}
            required
            sx={{ mb: 2 }}
            placeholder="e.g. Solved Leetcode Problem, Completed Algorithm Course"
          />
          
          <TextField
            fullWidth
            name="description"
            label="Description"
            value={formData.description}
            onChange={handleInputChange}
            multiline
            rows={4}
            sx={{ mb: 2 }}
            placeholder="Tell us more about what you did, what you learned, or any challenges you faced..."
          />
          
          <Box sx={{ mb: 2 }}>
            <Button
              variant="outlined"
              component="label"
              fullWidth
              sx={{ mb: 1 }}
            >
              {selectedImage ? 'Change Image' : 'Upload Image *'}
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
          
          <TextField
            fullWidth
            name="date"
            label="Date"
            type="date"
            value={formData.date}
            onChange={handleInputChange}
            required
            InputLabelProps={{
              shrink: true,
            }}
          />
          
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Group (Optional)</InputLabel>
            <Select
              name="group_id"
              value={formData.group_id || ''}
              onChange={handleSelectChange}
              label="Group (Optional)"
            >
              <MenuItem value="">
                <em>No Group (Personal Activity)</em>
              </MenuItem>
              {groups.map((group) => (
                <MenuItem key={group.id} value={group.id}>
                  {group.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
      </DialogContent>
      
      <DialogActions sx={{ p: 3 }}>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          color="secondary"
          disabled={loading || !formData.title || !formData.date || !selectedImage}
        >
          {loading ? 'Posting...' : 'Post Activity'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddActivityDialog;
