package monitoring

import (
	"context"
	"fmt"
	"log"
	"time"

	// "sync"

	"user-management/internal/domain/entities"

	"golang.org/x/sync/errgroup"
)

// NotificationService usa errgroup para manejo coordinado
type NotificationService struct {
	monitor   *OrderMonitor
	emailChan chan entities.Notification
	smsChan   chan entities.Notification
	pushChan  chan entities.Notification
}

func NewNotificationService(monitor *OrderMonitor) *NotificationService {
	return &NotificationService{
		monitor:   monitor,
		emailChan: make(chan entities.Notification, 50),
		smsChan:   make(chan entities.Notification, 20),
		pushChan:  make(chan entities.Notification, 10),
	}
}

// StartWorkers inicia múltiples workers con errgroup
func (ns *NotificationService) StartWorkers(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Worker de emails (3 instancias)
	for i := range 3 {
		workerID := i
		g.Go(func() error {
			return ns.emailWorker(ctx, workerID)
		})
	}

	// Worker de SMS (2 instancias)
	for i := range 2 {
		workerID := i
		g.Go(func() error {
			return ns.smsWorker(ctx, workerID)
		})
	}

	// Worker de push notifications
	g.Go(func() error {
		return ns.pushWorker(ctx)
	})

	// Worker de monitoreo
	g.Go(func() error {
		return ns.monitoringWorker(ctx)
	})

	// Esperar a que todos terminen o alguno falle
	return g.Wait()
}

func (ns *NotificationService) emailWorker(ctx context.Context, id int) error {
	log.Printf("Worker de email %d iniciado", id)

	for {
		select {
		case notification := <-ns.emailChan:
			// Simular envío
			time.Sleep(time.Millisecond * 100)
			log.Printf("Worker %d: Email enviado a %s",
				id, notification.Email)

			// Publicar evento al monitor
			ns.monitor.PublishEvent(OrderEvent{
				EventType: "email_sent",
				Data:      notification,
			})

		case <-ctx.Done():
			log.Printf("Worker de email %d terminado: %v",
				id, ctx.Err())
			return ctx.Err()
		}
	}
}

func (ns *NotificationService) smsWorker(ctx context.Context, id int) error {
	log.Printf("Worker de SMS %d iniciado", id)

	for {
		select {
		case notification := <-ns.smsChan:
			// Simular envío ocasionalmente fallido
			if time.Now().Unix()%10 == 0 { // 10% de fallos
				return fmt.Errorf("worker %d: error enviando SMS", id)
			}

			time.Sleep(time.Millisecond * 150)
			log.Printf("Worker %d: SMS enviado a %s",
				id, notification.Phone)

		case <-ctx.Done():
			log.Printf("Worker de SMS %d terminado", id)
			return nil
		}
	}
}

func (ns *NotificationService) pushWorker(ctx context.Context) error {
	log.Println("Worker de push notifications iniciado")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case notification := <-ns.pushChan:
			log.Printf("Push enviado: %s", notification.Message)

		case <-ticker.C:
			// Background task periódica
			log.Println("Push worker activo")

		case <-ctx.Done():
			log.Println("Worker de push terminado")
			return nil
		}
	}
}

func (ns *NotificationService) monitoringWorker(ctx context.Context) error {
	log.Println("Worker de monitoreo iniciado")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := ns.monitor.GetMetrics(time.Second)
			if err != nil {
				log.Printf("Error obteniendo métricas: %v", err)
			} else {
				log.Printf("Métricas actuales: %v", metrics)
			}

		case <-ctx.Done():
			log.Println("Worker de monitoreo terminado")
			return nil
		}
	}
}

// SendNotifications demuestra select con múltiples channels
func (ns *NotificationService) SendNotifications(
	notification entities.Notification,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Intentar enviar por diferentes canales con prioridad
	select {
	case ns.emailChan <- notification:
		return nil
	case ns.smsChan <- notification:
		return nil
	case ns.pushChan <- notification:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout enviando notificación: %w", ctx.Err())
	}
}
