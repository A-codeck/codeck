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
  FormControl,
  InputLabel,
  OutlinedInput,
  MenuItem,
  ListItemText,
  Checkbox,
  Chip,
} from '@mui/material';
import { Select, SelectChangeEvent } from '@mui/material';
import { ActivityCreateRequest, Group } from '../types/api';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';

interface AddActivityDialogProps {
  open: boolean;
  onClose: () => void;
  onActivityAdded: () => void;
  preselectedGroupId?: string;
}

const AddActivityDialog: React.FC<AddActivityDialogProps> = ({
  open,
  onClose,
  onActivityAdded,
  preselectedGroupId,
}) => {
  const { user } = useAuth();
  const [formData, setFormData] = useState<Omit<ActivityCreateRequest, 'image'>>({
    title: '',
    description: '',
    date: new Date().toISOString().split('T')[0], // Always current date
    group_ids: [], // Now an array
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
        creator_id: user.id,
        // Preselect group if provided
        group_ids: preselectedGroupId ? [preselectedGroupId] : []
      }));
      loadUserGroups();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user, preselectedGroupId]);

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
    
    // Skip date field as it's automatically set to current date
    if (name === 'date') return;
    
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

  const handleSelectChange = (event: SelectChangeEvent<string[]>) => {
    const value = event.target.value;
    setFormData({
      ...formData,
      group_ids: typeof value === 'string' ? value.split(',') : value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !selectedImage) return;

    try {
      setLoading(true);
      
      const activityData: ActivityCreateRequest = {
        ...formData,
        date: new Date().toISOString().split('T')[0], // Always use current date when posting
        image: selectedImage,
      };
      
      console.log('Creating activity with data:', {
        ...activityData,
        image: selectedImage.name // Don't log the full file object
      });
      
      const result = await apiService.createActivity(activityData);
      console.log('Activity created successfully:', result);
      
      // Reset form
      setFormData({
        title: '',
        description: '',
        date: new Date().toISOString().split('T')[0], // Always current date
        group_ids: [], // Reset to empty array
        creator_id: user.id, // Reset creator_id to current user
      });
      setSelectedImage(null);
      setImagePreview(null);
      
      console.log('Calling onActivityAdded callback');
      onActivityAdded();
      onClose();
    } catch (error: any) {
      console.error('Error creating activity:', error);
      console.error('Error details:', error.response?.data);
      alert('Failed to create activity. Please try again.');
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
          Registrar Sua Atividade
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          Compartilhe o que você conquistou hoje. Você deve selecionar pelo menos um grupo para compartilhar sua atividade.
        </Typography>
        <Typography variant="body2" color="primary" sx={{ fontWeight: 'medium' }}>
          📅 Será postado em: {new Date().toLocaleDateString('pt-BR', { 
            weekday: 'long', 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric' 
          })}
        </Typography>
      </DialogTitle>        <DialogContent>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            fullWidth
            name="title"
            label="Título da Atividade"
            value={formData.title}
            onChange={handleInputChange}
            required
            sx={{ mb: 2 }}
            placeholder="ex: Resolvi Problema no Leetcode, Completei Curso de Algoritmos"
          />
          
          <TextField
            fullWidth
            name="description"
            label="Descrição"
            value={formData.description}
            onChange={handleInputChange}
            multiline
            rows={4}
            sx={{ mb: 2 }}
            placeholder="Conte mais sobre o que você fez, o que aprendeu, ou quais desafios enfrentou..."
          />
          
          <Box sx={{ mb: 2 }}>
            <Button
              variant="outlined"
              component="label"
              fullWidth
              sx={{ mb: 1 }}
            >
              {selectedImage ? 'Alterar Imagem' : 'Fazer Upload da Imagem *'}
              <input
                type="file"
                hidden
                accept="image/*"
                onChange={handleImageChange}
              />
            </Button>
            {selectedImage && (
              <Typography variant="body2" color="text.secondary">
                Selecionado: {selectedImage.name}
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
          
          <FormControl fullWidth sx={{ mb: 2 }} required error={formData.group_ids.length === 0}>
            <InputLabel>Grupos *</InputLabel>
            <Select
              multiple
              value={formData.group_ids}
              onChange={handleSelectChange}
              input={<OutlinedInput label="Grupos *" />}
              disabled={groups.length === 0}
              renderValue={(selected) => (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {selected.map((value) => {
                    const group = groups.find(g => g.id === value);
                    return (
                      <Chip key={value} label={group?.name || value} size="small" />
                    );
                  })}
                </Box>
              )}
            >
              {groups.length === 0 ? (
                <MenuItem disabled>
                  <Typography color="text.secondary">
                    Nenhum grupo disponível. Por favor, crie ou entre em um grupo primeiro.
                  </Typography>
                </MenuItem>
              ) : (
                groups.map((group) => (
                  <MenuItem key={group.id} value={group.id}>
                    <Checkbox checked={formData.group_ids.indexOf(group.id) > -1} />
                    <ListItemText primary={group.name} />
                  </MenuItem>
                ))
              )}
            </Select>
            {groups.length === 0 ? (
              <Typography variant="caption" color="error" sx={{ mt: 0.5, ml: 1.5 }}>
                Você precisa fazer parte de pelo menos um grupo para postar atividades
              </Typography>
            ) : formData.group_ids.length === 0 ? (
              <Typography variant="caption" color="error" sx={{ mt: 0.5, ml: 1.5 }}>
                Por favor, selecione pelo menos um grupo
              </Typography>
            ) : preselectedGroupId && formData.group_ids.includes(preselectedGroupId) ? (
              <Typography variant="caption" color="primary" sx={{ mt: 0.5, ml: 1.5 }}>
                O grupo atual está pré-selecionado. Você pode selecionar grupos adicionais se desejar.
              </Typography>
            ) : null}
          </FormControl>
        </Box>
      </DialogContent>
      
      <DialogActions sx={{ p: 3 }}>
        <Button onClick={handleClose} disabled={loading}>
          Cancelar
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          color="secondary"
          disabled={loading || !formData.title || !selectedImage || formData.group_ids.length === 0 || groups.length === 0}
        >
          {loading ? 'Postando...' : 'Postar Atividade'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddActivityDialog;
