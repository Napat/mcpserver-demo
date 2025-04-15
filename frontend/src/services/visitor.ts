// API function สำหรับดึงข้อมูลจำนวนผู้เข้าชม
export const getVisitorCount = async (): Promise<number> => {
  try {
    const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const response = await fetch(`${API_URL}/api/visitors`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch visitor count');
    }

    const data = await response.json();
    return data.visitor_count;
  } catch (error) {
    console.error('Error fetching visitor count:', error);
    return 0;
  }
};

// API function สำหรับเพิ่มจำนวนผู้เข้าชม
export const incrementVisitorCount = async (): Promise<number> => {
  try {
    const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const response = await fetch(`${API_URL}/api/visitors`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to increment visitor count');
    }

    const data = await response.json();
    return data.visitor_count;
  } catch (error) {
    console.error('Error incrementing visitor count:', error);
    return 0;
  }
}; 